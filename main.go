package main

import (
	"K8sWatchDemo/pkg"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"sync"
	"time"
)

var nodeports = sync.Map{}
var eipports = sync.Map{}

func main() {
	go ApiWatchStart()
	go printUsedPorts(&nodeports)
	go printUsedPorts(&eipports)
	select {}

}

func printUsedPorts(ports *sync.Map) {
	for {
		time.Sleep(time.Second * 5)
		ports.Range(func(key, value interface{}) bool {
			fmt.Printf("%v  ", key)
			return true
		})
		fmt.Println()
	}

}

func ApiWatchStart() {
	helper := pkg.NewHelper()
	podListWatcher := cache.NewListWatchFromClient(
		helper.RESTClient(),
		"services",
		v1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewIndexerInformer(podListWatcher, &v1.Service{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			service := obj.(*v1.Service)
			fmt.Printf("service added: %s %s  \n", service.Namespace, service.Name)
			for _, v := range service.Spec.Ports {
				if v.NodePort > 0 {
					nodeports.Store(v.NodePort, v.NodePort)
				}

				if len(service.Spec.ExternalIPs) > 0 && v.Port > 0 {
					for _, eip := range service.Spec.ExternalIPs {
						eipports.Store(fmt.Sprintf("%s-%d", eip, v.Port), nil)
					}
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			service := obj.(*v1.Service)
			fmt.Printf("service deleted: %s  %s \n", service.Namespace, service.Name)
			for _, v := range service.Spec.Ports {
				nodeports.Delete(v.NodePort)

				if len(service.Spec.ExternalIPs) > 0 && v.Port > 0 {
					for _, eip := range service.Spec.ExternalIPs {
						eipports.Delete(fmt.Sprintf("%s-%d", eip, v.Port))
					}
				}
			}

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			os := oldObj.(*v1.Service)
			for _, v := range os.Spec.Ports {
				nodeports.Delete(v.NodePort)
				if len(os.Spec.ExternalIPs) > 0 && v.Port > 0 {
					for _, eip := range os.Spec.ExternalIPs {
						eipports.Delete(fmt.Sprintf("%s-%d", eip, v.Port))
					}
				}
			}

			service := newObj.(*v1.Service)
			fmt.Printf("service changed  %s   %s \n", service.Namespace, service.Name)
			for _, v := range service.Spec.Ports {
				nodeports.Store(v.NodePort, v.NodePort)

				if v.NodePort > 0 {
					nodeports.Store(v.NodePort, v.NodePort)
				}

				if len(service.Spec.ExternalIPs) > 0 && v.Port > 0 {
					for _, eip := range service.Spec.ExternalIPs {
						eipports.Store(fmt.Sprintf("%s-%d", eip, v.Port), nil)
					}
				}

			}
		},
	}, cache.Indexers{})
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
