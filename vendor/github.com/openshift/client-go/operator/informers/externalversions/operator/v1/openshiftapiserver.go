// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	context "context"
	time "time"

	apioperatorv1 "github.com/openshift/api/operator/v1"
	versioned "github.com/openshift/client-go/operator/clientset/versioned"
	internalinterfaces "github.com/openshift/client-go/operator/informers/externalversions/internalinterfaces"
	operatorv1 "github.com/openshift/client-go/operator/listers/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// OpenShiftAPIServerInformer provides access to a shared informer and lister for
// OpenShiftAPIServers.
type OpenShiftAPIServerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() operatorv1.OpenShiftAPIServerLister
}

type openShiftAPIServerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewOpenShiftAPIServerInformer constructs a new informer for OpenShiftAPIServer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewOpenShiftAPIServerInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredOpenShiftAPIServerInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredOpenShiftAPIServerInformer constructs a new informer for OpenShiftAPIServer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredOpenShiftAPIServerInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OperatorV1().OpenShiftAPIServers().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OperatorV1().OpenShiftAPIServers().Watch(context.TODO(), options)
			},
		},
		&apioperatorv1.OpenShiftAPIServer{},
		resyncPeriod,
		indexers,
	)
}

func (f *openShiftAPIServerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredOpenShiftAPIServerInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *openShiftAPIServerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apioperatorv1.OpenShiftAPIServer{}, f.defaultInformer)
}

func (f *openShiftAPIServerInformer) Lister() operatorv1.OpenShiftAPIServerLister {
	return operatorv1.NewOpenShiftAPIServerLister(f.Informer().GetIndexer())
}
