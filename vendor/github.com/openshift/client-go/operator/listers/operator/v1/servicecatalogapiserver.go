// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
)

// ServiceCatalogAPIServerLister helps list ServiceCatalogAPIServers.
// All objects returned here must be treated as read-only.
type ServiceCatalogAPIServerLister interface {
	// List lists all ServiceCatalogAPIServers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*operatorv1.ServiceCatalogAPIServer, err error)
	// Get retrieves the ServiceCatalogAPIServer from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*operatorv1.ServiceCatalogAPIServer, error)
	ServiceCatalogAPIServerListerExpansion
}

// serviceCatalogAPIServerLister implements the ServiceCatalogAPIServerLister interface.
type serviceCatalogAPIServerLister struct {
	listers.ResourceIndexer[*operatorv1.ServiceCatalogAPIServer]
}

// NewServiceCatalogAPIServerLister returns a new ServiceCatalogAPIServerLister.
func NewServiceCatalogAPIServerLister(indexer cache.Indexer) ServiceCatalogAPIServerLister {
	return &serviceCatalogAPIServerLister{listers.New[*operatorv1.ServiceCatalogAPIServer](indexer, operatorv1.Resource("servicecatalogapiserver"))}
}
