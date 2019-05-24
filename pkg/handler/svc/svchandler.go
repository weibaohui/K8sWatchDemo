package svc

import (
	"K8sWatchDemo/pkg/event"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type ServiceHandler struct {
	logger *logrus.Entry
}

func (h *ServiceHandler) Init() {
	h.logger = logrus.WithField("handler", "ServiceHandler")
}
func (h *ServiceHandler) ObjectCreated(obj interface{}) {
	svc := obj.(*v1.Service)
	h.logger.Infoln("add", svc.Name)
}

func (h *ServiceHandler) ObjectDeleted(event event.InformerEvent) {
	h.logger.Infoln("delete", event)
}

func (h *ServiceHandler) ObjectUpdated(oldObj interface{}, event event.InformerEvent) {
	svc := oldObj.(*v1.Service)
	h.logger.Infoln("update", svc.Name)
}
