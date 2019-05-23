package handler

import (
	"K8sWatchDemo/event"
	"fmt"
)

type Handler interface {
	ObjectCreated(obj interface{})
	ObjectDeleted(event event.InformerEvent)
	ObjectUpdated(oldObj interface{}, event event.InformerEvent)
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

func (d *DefaultHandler) ObjectCreated(obj interface{}) {
	fmt.Println("DefaultHandler ObjectCreated ", obj)

}

func (d *DefaultHandler) ObjectDeleted(event event.InformerEvent) {
	fmt.Println("DefaultHandler ObjectDeleted ", event)

}

func (d *DefaultHandler) ObjectUpdated(oldObj interface{}, event event.InformerEvent) {
	fmt.Println("DefaultHandler ObjectUpdated ", event)

}
