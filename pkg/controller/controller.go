package controller

import (
	"K8sWatchDemo/pkg/config"
	"K8sWatchDemo/pkg/event"
	"K8sWatchDemo/pkg/handler"
	"K8sWatchDemo/pkg/utils"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const maxRetries = 5

type Controller struct {
	logger       *logrus.Entry
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler handler.Handler
}

func Start(conf *config.Config) {
	var kubeClient = utils.NewHelper().GetKubeClient()

	for _, v := range conf.Handlers {
		v.Init()
	}

	if conf.Resource.Pod {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Pods(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Pods(conf.Namespace).Watch(options)
				},
			},
			&api_v1.Pod{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["po"], informer, "pod")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)

		informerHeadless := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Pods(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Pods(conf.Namespace).Watch(options)
				},
			},
			&api_v1.Pod{},
			0,
			cache.Indexers{},
		)

		hc := newResourceController(kubeClient, conf.Handlers["headless-po"], informerHeadless, "pod")
		stopChHc := make(chan struct{})
		defer close(stopChHc)

		go hc.Run(stopCh)

	}

	if conf.Resource.DaemonSet {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.ExtensionsV1beta1().DaemonSets(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.ExtensionsV1beta1().DaemonSets(conf.Namespace).Watch(options)
				},
			},
			&ext_v1beta1.DaemonSet{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["ds"], informer, "daemonset")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.ReplicaSet {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.ExtensionsV1beta1().ReplicaSets(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.ExtensionsV1beta1().ReplicaSets(conf.Namespace).Watch(options)
				},
			},
			&ext_v1beta1.ReplicaSet{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["rs"], informer, "replicaset")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Services {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Services(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Services(conf.Namespace).Watch(options)
				},
			},
			&api_v1.Service{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["svc"], informer, "service")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Deployment {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1beta1().Deployments(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1beta1().Deployments(conf.Namespace).Watch(options)
				},
			},
			&apps_v1beta1.Deployment{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["deploy"], informer, "deployment")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Namespace {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Namespaces().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Namespaces().Watch(options)
				},
			},
			&api_v1.Namespace{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["ns"], informer, "namespace")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.ReplicationController {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ReplicationControllers(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ReplicationControllers(conf.Namespace).Watch(options)
				},
			},
			&api_v1.ReplicationController{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["rc"], informer, "replication controller")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Job {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.BatchV1().Jobs(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.BatchV1().Jobs(conf.Namespace).Watch(options)
				},
			},
			&batch_v1.Job{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["job"], informer, "job")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.PersistentVolume {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().PersistentVolumes().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().PersistentVolumes().Watch(options)
				},
			},
			&api_v1.PersistentVolume{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["pv"], informer, "persistent volume")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Secret {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Secrets(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Secrets(conf.Namespace).Watch(options)
				},
			},
			&api_v1.Secret{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["secret"], informer, "secret")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.ConfigMap {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ConfigMaps(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ConfigMaps(conf.Namespace).Watch(options)
				},
			},
			&api_v1.ConfigMap{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["cm"], informer, "configmap")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Ingress {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.ExtensionsV1beta1().Ingresses(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.ExtensionsV1beta1().Ingresses(conf.Namespace).Watch(options)
				},
			},
			&ext_v1beta1.Ingress{},
			0,
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, conf.Handlers["ing"], informer, "ingress")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func GetNamespace(key string) (namespace, name string) {
	if strings.Contains(key, "/") {
		names := strings.SplitN(key, "/", 2)
		namespace = names[0]
		name = names[1]
		return
	}
	return "", ""
}
func newResourceController(client kubernetes.Interface, eventHandler handler.Handler, informer cache.SharedIndexInformer, resourceType string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var eve event.InformerEvent
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			eve.Key, err = cache.MetaNamespaceKeyFunc(obj)
			eve.EventType = "create"
			eve.ResourceType = resourceType
			eve.Namespace, eve.Name = GetNamespace(eve.Key)
			logrus.WithField("pkg", "k8swatch-"+resourceType).Infof("新增 %v: %s", resourceType, eve.Key)
			if err == nil {
				queue.Add(eve)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			eve.Key, err = cache.MetaNamespaceKeyFunc(old)
			eve.EventType = "update"
			eve.ResourceType = resourceType
			eve.Namespace, eve.Name = GetNamespace(eve.Key)
			logrus.WithField("pkg", "k8swatch-"+resourceType).Infof("更新 %v: %s", resourceType, eve.Key)
			if err == nil {
				queue.Add(eve)
			}
		},
		DeleteFunc: func(obj interface{}) {
			eve.Key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			eve.EventType = "delete"
			eve.ResourceType = resourceType
			eve.Namespace, eve.Name = GetNamespace(eve.Key)
			logrus.WithField("pkg", "k8swatch-"+resourceType).Infof("删除 %v: %s", resourceType, eve.Key)
			if err == nil {
				queue.Add(eve)
			}
		},
	})

	return &Controller{
		logger:       logrus.WithField("pkg", "k8swatch-"+resourceType),
		clientset:    client,
		informer:     informer,
		queue:        queue,
		eventHandler: eventHandler,
	}
}

// Run starts the k8swatch controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("监控程序启动")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("等待缓存同步超时"))
		return
	}

	c.logger.Info("监控程序 已完成同步，开始工作")

	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	eve, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(eve)
	err := c.processItem(eve.(event.InformerEvent))
	if err == nil {
		c.queue.Forget(eve)
	} else if c.queue.NumRequeues(eve) < maxRetries {
		c.logger.Errorf("处理 %s 出错，进行重试: %v", eve.(event.InformerEvent).Key, err)
		c.queue.AddRateLimited(eve)
	} else {
		c.logger.Errorf("处理 %s 出错，丢弃: %v", eve.(event.InformerEvent).Key, err)
		c.queue.Forget(eve)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(e event.InformerEvent) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(e.Key)
	if err != nil {
		return fmt.Errorf("获取 key %s 出错: %v", e.Key, err)
	}
	switch e.EventType {
	case "create":
		c.eventHandler.ObjectCreated(obj)
	case "update":
		c.eventHandler.ObjectUpdated(obj, e)
		return nil
	case "delete":
		c.eventHandler.ObjectDeleted(e)
		return nil
	}
	return nil
}
