package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"os"
	"path/filepath"
)

func main() {
	cli := getClient()
	// checkPod("dubbo", cli)

	podListWatcher := cache.NewListWatchFromClient(
		cli.CoreV1().RESTClient(),
		"pods",
		v1.NamespaceDefault,
		fields.Everything())

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
			fmt.Println("ADD")
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
			fmt.Println("UPDATE")
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
			fmt.Println("DELETE")
		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer)

	// We can now warm up the cache for initial synchronization.
	// Let's suppose that we knew about a pod "mypod" on our last run, therefore add it to the cache.
	// If this pod is not there anymore, the controller will be notified about the removal after the
	// cache has synchronized.
	indexer.Add(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: v1.NamespaceDefault,
			Labels: map[string]string{
				"app": "dubbo",
			},
		},
	})

	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// Wait forever
	select {}
}

func getClient() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		s := filepath.Join(home, ".kube", "config")
		kubeconfig = flag.String("kubeconfig", s, "kubeconfig存放位置")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig存放位置")
	}
	flag.Parse()
	fmt.Println(*kubeconfig)
	config, e := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if e != nil {
		panic(e.Error())
	}
	clientset, e := kubernetes.NewForConfig(config)
	if e != nil {
		panic(e.Error())
	}
	return clientset
}

func homeDir() string {
	if s := os.Getenv("HOME"); s != "" {
		return s
	}
	return os.Getenv("USERPROFILE")
}
