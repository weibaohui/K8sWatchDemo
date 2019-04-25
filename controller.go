package main

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	cli      *kubernetes.Clientset
}

func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller, cli *kubernetes.Clientset) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
		cli:      cli,
	}

}

func (c *Controller) processNextItem() bool {
	act, quit := c.queue.Get()
	if quit {
		return false
	}

	defer c.queue.Done(act)

	err := c.syncToStdout(act.(Action), c.cli)

	c.handleErr(err, act)
	return true
}

func (c *Controller) syncToStdout(act Action, cli *kubernetes.Clientset) error {
	_, exists, err := c.indexer.GetByKey(act.PodName)
	if err != nil {
		return err
	}

	fmt.Println("收到消息", act.ActionName, act.PodName)
	namespace, podName := getPodName(act.PodName)
	fmt.Println("所属namespace", namespace)
	if !isTarget(podName) {
		return nil
	}

	switch act.ActionName {
	case ADD:
		addPodProcess(cli, podName)
	case DELETE:
		if !exists {
			deletePodProcess(cli, podName)
		}
	case UPDATE:
	}

	return nil
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	klog.Infof("Dropping pod %q out of the queue: %v", key, err)
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}
