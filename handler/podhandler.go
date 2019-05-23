package handler

import (
	"K8sWatchDemo/event"
	"fmt"
	v1 "k8s.io/api/core/v1"
)

type PodHandler struct {
}

func (d *PodHandler) ObjectCreated(obj interface{}) {
	pod := obj.(*v1.Pod)
	fmt.Println("PodHandler ObjectCreated ", pod.Name)

}

func (d *PodHandler) ObjectDeleted(event event.InformerEvent) {
	fmt.Println("PodHandler ObjectDeleted ", event)

}

func (d *PodHandler) ObjectUpdated(oldObj interface{}, event event.InformerEvent) {
	pod := oldObj.(*v1.Pod)
	fmt.Println("PodHandler ObjectUpdated ", pod.Name)

}
