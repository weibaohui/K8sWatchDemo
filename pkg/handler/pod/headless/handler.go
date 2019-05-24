package headless

import (
	"K8sWatchDemo/pkg/event"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type HeadlessPodHandler struct {
	logger *logrus.Entry
}

func (h *HeadlessPodHandler) Init() {
	h.logger = logrus.WithField("handler", "HeadlessPodHandler")
}
func (h *HeadlessPodHandler) ObjectCreated(obj interface{}) {
	pod := obj.(*v1.Pod)
	h.logger.Infoln("add", pod.Name)

}

func (h *HeadlessPodHandler) ObjectDeleted(event event.InformerEvent) {
	h.logger.Infoln("delete", event)

}

func (h *HeadlessPodHandler) ObjectUpdated(oldObj interface{}, event event.InformerEvent) {
	pod := oldObj.(*v1.Pod)
	h.logger.Infoln("update", pod.Name)

}
