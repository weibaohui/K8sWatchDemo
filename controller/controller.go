package controller

import (
	"K8sWatchDemo/pkg"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"

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

var serverStartTime time.Time

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

// Controller object
type Controller struct {
	logger       *logrus.Entry
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler Handler
}

func Start(conf *Config, eventHandler Handler) {
	var kubeClient = pkg.NewHelper().GetKubeClient()

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

		c := newResourceController(kubeClient, eventHandler, informer, "pod")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
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

		c := newResourceController(kubeClient, eventHandler, informer, "daemonset")
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

		c := newResourceController(kubeClient, eventHandler, informer, "replicaset")
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

		c := newResourceController(kubeClient, eventHandler, informer, "service")
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

		c := newResourceController(kubeClient, eventHandler, informer, "deployment")
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

		c := newResourceController(kubeClient, eventHandler, informer, "namespace")
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

		c := newResourceController(kubeClient, eventHandler, informer, "replication controller")
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

		c := newResourceController(kubeClient, eventHandler, informer, "job")
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

		c := newResourceController(kubeClient, eventHandler, informer, "persistent volume")
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

		c := newResourceController(kubeClient, eventHandler, informer, "secret")
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

		c := newResourceController(kubeClient, eventHandler, informer, "configmap")
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

		c := newResourceController(kubeClient, eventHandler, informer, "ingress")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func GetNamespace(key string) (namespace string) {
	if strings.Contains(key, "/") {
		names := strings.SplitN(key, "/", 2)
		namespace = names[0]
		return
	}
	return ""
}
func newResourceController(client kubernetes.Interface, eventHandler Handler, informer cache.SharedIndexInformer, resourceType string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent Event
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(obj)
			newEvent.eventType = "create"
			newEvent.resourceType = resourceType
			newEvent.namespace = GetNamespace(newEvent.key)
			logrus.WithField("pkg", "k8swatch-"+resourceType).Infof("Processing add to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
			newEvent.eventType = "update"
			newEvent.resourceType = resourceType
			newEvent.namespace = GetNamespace(newEvent.key)
			logrus.WithField("pkg", "k8swatch-"+resourceType).Infof("Processing update to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
		DeleteFunc: func(obj interface{}) {
			newEvent.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			newEvent.eventType = "delete"
			newEvent.resourceType = resourceType
			newEvent.namespace = GetNamespace(newEvent.key)
			logrus.WithField("pkg", "k8swatch-"+resourceType).Infof("Processing delete to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
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

	c.logger.Info("Starting k8swatch controller")
	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	c.logger.Info("k8swatch controller synced and ready")

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
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(Event))
	if err == nil {
		// No error, reset the rate limit counters
		c.queue.Forget(newEvent)
	} else if c.queue.NumRequeues(newEvent) < maxRetries {
		c.logger.Errorf("Error processing %s (will retry): %v", newEvent.(Event).key, err)
		c.queue.AddRateLimited(newEvent)
	} else {
		// err != nil and too many retries
		c.logger.Errorf("Error processing %s (giving up): %v", newEvent.(Event).key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(newEvent Event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.key)
	if err != nil {
		return fmt.Errorf("获取 key %s 出错: %v", newEvent.key, err)
	}

	// process events based on its type
	switch newEvent.eventType {
	case "create":
		c.eventHandler.ObjectCreated(obj)
	case "update":
		kbEvent := KubeEvent{
			Kind:      newEvent.resourceType,
			Name:      newEvent.key,
			Namespace: newEvent.namespace,
		}
		c.eventHandler.ObjectUpdated(obj, kbEvent)
		return nil
	case "delete":
		kbEvent := KubeEvent{
			Kind:      newEvent.resourceType,
			Name:      newEvent.key,
			Namespace: newEvent.namespace,
		}
		c.eventHandler.ObjectDeleted(kbEvent)
		return nil
	}
	return nil
}
