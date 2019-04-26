package watcher

import (
	"K8sWatchDemo/pkg"
	"K8sWatchDemo/utils"
	"fmt"
	v1 "k8s.io/api/core/v1"
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
	helper   *pkg.Helper
}

func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller, helper *pkg.Helper) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
		helper:   helper,
	}

}

func (c *Controller) processNextItem() bool {
	act, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(act)
	err := c.processEvent(act.(pkg.Action))
	c.handleErr(err, act)
	return true
}

func (c *Controller) processEvent(act pkg.Action) error {
	obj, exists, err := c.indexer.GetByKey(act.PodNameNs)
	if err != nil {
		return err
	}

	namespace, podName := utils.GetPodName(act.PodNameNs)
	if !isTarget(namespace, podName, isTargetByPodName) {
		return nil
	}

	switch act.ActionName {
	case ADD:
		c.helper.AddProcess(obj.(*v1.Pod))
	case DELETE:
		if !exists {
			c.helper.DeleteProcess(act.PodNameNs)
		}
	case UPDATE:
		c.helper.UpdateProcess(act.PodNameNs)
	}

	return nil
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		fmt.Printf("Error syncing pod %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
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
