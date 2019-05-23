package handler

import (
	"K8sWatchDemo/controller"
	"fmt"
)

type Handler interface {
	Init(c *controller.Config) error
	ObjectCreated(obj interface{})
	ObjectDeleted(event ResourceType)
	ObjectUpdated(oldObj interface{}, event ResourceType)
}

var Map = map[string]Handler{
	"default": &DefaultHandler{},
	"deploy":  &DefaultHandler{},
	"rc":      &DefaultHandler{},
	"rs":      &DefaultHandler{},
	"ds":      &DefaultHandler{},
	"svc":     &DefaultHandler{},
	"po":      &PodHandler{},
	"job":     &DefaultHandler{},
	"pv":      &DefaultHandler{},
	"cm":      &DefaultHandler{},
	"ns":      &DefaultHandler{},
	"ing":     &DefaultHandler{},
	"secret":  &DefaultHandler{},
}

type DefaultHandler struct {
}

func (d *DefaultHandler) Init(c *controller.Config) error {
	fmt.Println("DefaultHandler init ")

	return nil
}

func (d *DefaultHandler) ObjectCreated(obj interface{}) {
	fmt.Println("DefaultHandler ObjectCreated ", obj)

}

func (d *DefaultHandler) ObjectDeleted(event ResourceType) {
	fmt.Println("DefaultHandler ObjectDeleted ", event)

}

func (d *DefaultHandler) ObjectUpdated(oldObj interface{}, event ResourceType) {
	fmt.Println("DefaultHandler ObjectUpdated ", event)

}
