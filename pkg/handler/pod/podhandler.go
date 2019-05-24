package pod

import (
	"K8sWatchDemo/pkg/event"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type PodHandler struct {
	logger *logrus.Entry
}

func (h *PodHandler) Init() {
	h.logger = logrus.WithField("handler", "PodHandler")
}
func (h *PodHandler) ObjectCreated(obj interface{}) {
	pod := obj.(*v1.Pod)
	h.logger.Infoln("add", pod.Name)

}

func (h *PodHandler) ObjectDeleted(event event.InformerEvent) {
	h.logger.Infoln("delete", event)
}

func (h *PodHandler) ObjectUpdated(oldObj interface{}, event event.InformerEvent) {
	pod := oldObj.(*v1.Pod)
	h.logger.Infoln("update", pod.Name)
}
