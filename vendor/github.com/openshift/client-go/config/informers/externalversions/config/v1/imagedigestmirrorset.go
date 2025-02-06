// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	context "context"
	time "time"

	apiconfigv1 "github.com/openshift/api/config/v1"
	versioned "github.com/openshift/client-go/config/clientset/versioned"
	internalinterfaces "github.com/openshift/client-go/config/informers/externalversions/internalinterfaces"
	configv1 "github.com/openshift/client-go/config/listers/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ImageDigestMirrorSetInformer provides access to a shared informer and lister for
// ImageDigestMirrorSets.
type ImageDigestMirrorSetInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() configv1.ImageDigestMirrorSetLister
}

type imageDigestMirrorSetInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewImageDigestMirrorSetInformer constructs a new informer for ImageDigestMirrorSet type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewImageDigestMirrorSetInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredImageDigestMirrorSetInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredImageDigestMirrorSetInformer constructs a new informer for ImageDigestMirrorSet type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredImageDigestMirrorSetInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().ImageDigestMirrorSets().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().ImageDigestMirrorSets().Watch(context.TODO(), options)
			},
		},
		&apiconfigv1.ImageDigestMirrorSet{},
		resyncPeriod,
		indexers,
	)
}

func (f *imageDigestMirrorSetInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredImageDigestMirrorSetInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *imageDigestMirrorSetInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiconfigv1.ImageDigestMirrorSet{}, f.defaultInformer)
}

func (f *imageDigestMirrorSetInformer) Lister() configv1.ImageDigestMirrorSetLister {
	return configv1.NewImageDigestMirrorSetLister(f.Informer().GetIndexer())
}
