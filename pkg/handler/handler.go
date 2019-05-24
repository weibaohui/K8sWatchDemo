package handler

import (
	"K8sWatchDemo/pkg/event"
	"K8sWatchDemo/pkg/handler/cm"
	"K8sWatchDemo/pkg/handler/pod"
	"K8sWatchDemo/pkg/handler/pod/headless"
	"K8sWatchDemo/pkg/handler/svc"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	Init()
	ObjectCreated(obj interface{})
	ObjectDeleted(event event.InformerEvent)
	ObjectUpdated(oldObj interface{}, event event.InformerEvent)
}

var Map = map[string]Handler{
	"default":     &DefaultHandler{},
	"deploy":      &DefaultHandler{},
	"rc":          &DefaultHandler{},
	"rs":          &DefaultHandler{},
	"ds":          &DefaultHandler{},
	"svc":         &svc.ServiceHandler{},
	"po":          &pod.PodHandler{},
	"headless-po": &headless.HeadlessPodHandler{},
	"job":         &DefaultHandler{},
	"pv":          &DefaultHandler{},
	"cm":          &cm.ConfigMapHandler{},
	"ns":          &DefaultHandler{},
	"ing":         &DefaultHandler{},
	"secret":      &DefaultHandler{},
}

type DefaultHandler struct {
	logger *logrus.Entry
}

func (h *DefaultHandler) Init() {
	h.logger = logrus.WithField("handler", "DefaultHandler")
}
func (h *DefaultHandler) ObjectCreated(obj interface{}) {
	h.logger.Infoln("add", obj)
}

func (h *DefaultHandler) ObjectDeleted(event event.InformerEvent) {
	h.logger.Infoln("delete", event)

}

func (h *DefaultHandler) ObjectUpdated(oldObj interface{}, event event.InformerEvent) {
	h.logger.Infoln("update", event)

}
