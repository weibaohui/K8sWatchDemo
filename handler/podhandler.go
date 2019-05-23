package handler

import (
	"K8sWatchDemo/controller"
	"fmt"
)

type PodHandler struct {
}

func (d *PodHandler) Init(c *controller.Config) error {
	fmt.Println("PodHandler init ")

	return nil
}

func (d *PodHandler) ObjectCreated(obj interface{}) {
	fmt.Println("PodHandler ObjectCreated ", obj)

}

func (d *PodHandler) ObjectDeleted(event ResourceType) {
	fmt.Println("PodHandler ObjectDeleted ", event)

}

func (d *PodHandler) ObjectUpdated(oldObj interface{}, event ResourceType) {
	fmt.Println("PodHandler ObjectUpdated ", event)

}
