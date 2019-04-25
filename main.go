package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"os"
	"path/filepath"
)

func main() {
	cli := getClient()

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
				queue.Add(Action{
					PodName:    key,
					ActionName: ADD,
				})
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(Action{
					PodName:    key,
					ActionName: UPDATE,
				})
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(Action{
					PodName:    key,
					ActionName: DELETE,
				})
			}

		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer, cli)

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// Wait forever
	select {}
}

func getClient() *kubernetes.Clientset {
	var kubeconfig *string
	var inCluster *bool
	if home := homeDir(); home != "" {
		s := filepath.Join(home, ".kube", "config")
		kubeconfig = flag.String("kubeconfig", s, "kubeconfig存放位置")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig存放位置")
	}
	inCluster = flag.Bool("in", false, "是否在集群内")
	flag.Parse()
	var config *rest.Config
	var err error
	if *inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

	} else {
		fmt.Println(*kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

	}
	cli, e := kubernetes.NewForConfig(config)
	if e != nil {
		panic(e.Error())
	}
	return cli

}

func homeDir() string {
	if s := os.Getenv("HOME"); s != "" {
		return s
	}
	return os.Getenv("USERPROFILE")
}
