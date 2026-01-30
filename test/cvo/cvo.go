package cvo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"

	"github.com/openshift/cluster-version-operator/test/oc"
	ocapi "github.com/openshift/cluster-version-operator/test/oc/api"
	"github.com/openshift/cluster-version-operator/test/util"
)

var logger = g.GinkgoLogr.WithName("cluster-version-operator-tests")

const cvoNamespace = "openshift-cluster-version"

var _ = g.Describe(`[Jira:"Cluster Version Operator"] cluster-version-operator-tests`, func() {
	g.It("should support passing tests", func() {
		o.Expect(true).To(o.BeTrue())
	})

	g.It("can use oc to get the version information", func() {
		ocClient, err := oc.NewOC(logger)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(ocClient).NotTo(o.BeNil())

		output, err := ocClient.Version(ocapi.VersionOptions{Client: true})
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(output).To(o.ContainSubstring("Client Version:"))
	})
})

// CVO tests which need access the live cluster will be placed here
var _ = g.Describe(`[Jira:"Cluster Version Operator"] cluster-version-operator`, func() {
	var (
		restCfg    *rest.Config
		kubeClient kubernetes.Interface
	)

	g.BeforeEach(func() {
		var err error
		// Respects KUBECONFIG env var
		restCfg, err = util.GetRestConfig()
		o.Expect(err).NotTo(o.HaveOccurred(), "Failed to load Kubernetes configuration. Please ensure KUBECONFIG environment variable is set.")

		kubeClient, err = util.GetKubeClient(restCfg)
		o.Expect(err).NotTo(o.HaveOccurred(), "Failed to create Kubernetes client")
	})

	// Migrated from case NonHyperShiftHOST-Author:jiajliu-Low-46922-check runlevel and scc in cvo ns
	// Refer to https://github.com/openshift/openshift-tests-private/blob/40374cf20946ff03c88712839a5626af2c88ab31/test/extended/ota/cvo/cvo.go#L1081
	g.It("should have correct runlevel and scc", func() {
		ctx := context.Background()
		err := util.SkipIfHypershift(ctx, restCfg)
		o.Expect(err).NotTo(o.HaveOccurred(), "Failed to determine if cluster is HyperShift")
		err = util.SkipIfMicroshift(ctx, restCfg)
		o.Expect(err).NotTo(o.HaveOccurred(), "Failed to determine if cluster is MicroShift")

		g.By("Checking that the 'openshift.io/run-level' label exists on the namespace and has the empty value")
		ns, err := kubeClient.CoreV1().Namespaces().Get(ctx, cvoNamespace, metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred(), "Failed to get namespace %s", cvoNamespace)
		runLevel, exists := ns.Labels["openshift.io/run-level"]
		o.Expect(exists).To(o.BeTrue(), "The 'openshift.io/run-level' label on namespace %s does not exist", cvoNamespace)
		o.Expect(runLevel).To(o.BeEmpty(), "Expected the 'openshift.io/run-level' label value on namespace %s has the empty value, but got %s", cvoNamespace, runLevel)

		g.By("Checking that the annotation 'openshift.io/scc annotation' on the CVO pod has the value hostaccess")
		podList, err := kubeClient.CoreV1().Pods(cvoNamespace).List(ctx, metav1.ListOptions{
			LabelSelector: "k8s-app=cluster-version-operator",
			FieldSelector: "status.phase=Running",
		})
		o.Expect(err).NotTo(o.HaveOccurred(), "Failed to list running CVO pods")
		o.Expect(podList.Items).To(o.HaveLen(1), "Expected exactly one running CVO pod, but found: %d", len(podList.Items))

		cvoPod := podList.Items[0]
		sccAnnotation := cvoPod.Annotations["openshift.io/scc"]
		o.Expect(sccAnnotation).To(o.Equal("hostaccess"), "Expected the annotation 'openshift.io/scc annotation' on pod %s to have the value 'hostaccess', but got %s", cvoPod.Name, sccAnnotation)
	})

	g.It(`should not install resources annotated with release.openshift.io/delete=true`, g.Label("Conformance", "High", "42543"), func() {
		// Initialize the ocapi.OC instance
		g.By("Setup ocapi.OC")
		err := os.Setenv("OC_CLI_TIMEOUT", "90s")
		o.Expect(err).NotTo(o.HaveOccurred(), "Setup environment variable OC_CLI_TIMEOUT failed")
		ocClient, err := oc.NewOC(logger)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(ocClient).NotTo(o.BeNil())
		defer func() {
			err = os.Unsetenv("OC_CLI_TIMEOUT")
			o.Expect(err).NotTo(o.HaveOccurred(), "Unset environment variable OC_CLI_TIMEOUT failed")
		}()

		g.By("Extract manifests")
		annotation := "release.openshift.io/delete"
		manifestDir := ocapi.ReleaseExtractOptions{To: "/tmp/OTA-42543-manifest"}
		logger.Info(fmt.Sprintf("Extract manifests to: %s", manifestDir.To))
		defer func() { _ = os.RemoveAll(manifestDir.To) }()
		err = ocClient.AdmReleaseExtract(manifestDir)
		o.Expect(err).NotTo(o.HaveOccurred(), "extract manifests failed")

		entries, err := os.ReadDir(manifestDir.To)
		o.Expect(err).NotTo(o.HaveOccurred())
		g.By("Start to iterate all manifests")
		var closeFilePass = true
		for _, entry := range entries {
			nameLower := strings.ToLower(entry.Name())
			if strings.Contains(nameLower, "cleanup") {
				logger.Info(fmt.Sprintf("Skipping file %s because it matches cleanup filter", entry.Name()))
				continue
			}
			filePath := filepath.Join(manifestDir.To, entry.Name())
			file, err := os.Open(filePath)
			o.Expect(err).NotTo(o.HaveOccurred())
			defer func() {
				if !closeFilePass {
					// Close the file again
					if err = file.Close(); err != nil {
						o.Expect(err).NotTo(o.HaveOccurred(), "close file failed")
					}
				}
			}()
			decoder := yamlv3.NewDecoder(file)
			for {
				var doc map[string]interface{}
				if err := decoder.Decode(&doc); err != nil {
					if err == io.EOF {
						break
					}
					continue
				}
				meta, _ := doc["metadata"].(map[string]interface{})
				ann, _ := meta["annotations"].(map[string]interface{})
				if ann == nil || ann[annotation] != "true" {
					continue
				}
				kind, _ := doc["kind"].(string)
				name, _ := meta["name"].(string)
				namespace, _ := meta["namespace"].(string)
				args := []string{"get", kind, name}
				if namespace != "" {
					args = append(args, "-n", namespace)
				}
				_, err := ocClient.Run(args...)
				o.Expect(err).To(o.HaveOccurred(), "The deleted manifest should not be installed, but actually installed")
			}
			// close each file
			err = file.Close()
			if err != nil {
				closeFilePass = false
				o.Expect(err).NotTo(o.HaveOccurred(), "close file failed")
			}
		}
	})
})
