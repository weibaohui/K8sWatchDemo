package main

import (
	"K8sWatchDemo/pkg"
	"K8sWatchDemo/watcher"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func main() {

	helper := pkg.NewHelper()

	podListWatcher := cache.NewListWatchFromClient(
		helper.RESTClient(),
		"pods",
		v1.NamespaceDefault,
		fields.Everything())

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			podNameNs, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(pkg.Action{
					PodNameNs:  podNameNs,
					ActionName: watcher.ADD,
				})
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			podNameNs, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(pkg.Action{
					PodNameNs:  podNameNs,
					ActionName: watcher.UPDATE,
				})
			}
		},
		DeleteFunc: func(obj interface{}) {
			podNameNs, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(pkg.Action{
					PodNameNs:  podNameNs,
					ActionName: watcher.DELETE,
				})
			}

		},
	}, cache.Indexers{})

	controller := watcher.NewController(queue, indexer, informer, helper)

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// Wait forever
	select {}
}
