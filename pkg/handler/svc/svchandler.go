package svc

import (
	"K8sWatchDemo/pkg/cluster"
	"K8sWatchDemo/pkg/event"
	"K8sWatchDemo/pkg/utils"
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
	addPortConfigs(svc)
	h.logger.Infoln("add", svc.Name)
}
func (h *ServiceHandler) ObjectDeleted(event event.InformerEvent) {
	cluster.GetClusterConfig().DeleteSvc(event.Namespace, event.Name)
	h.logger.Infoln("delete", event)
}

func (h *ServiceHandler) ObjectUpdated(oldObj interface{}, event event.InformerEvent) {
	svc := oldObj.(*v1.Service)

	cluster.GetClusterConfig().DeleteSvc(svc.Namespace, svc.Name)
	addPortConfigs(svc)
	h.logger.Infoln("update", svc.Name)
}

func addPortConfigs(svc *v1.Service) {
	for _, p := range svc.Spec.Ports {
		if p.NodePort > 0 {
			//存在NodePort

			for _, nip := range utils.NodeIPs() {
				cluster.GetClusterConfig().Add(&cluster.IpPortConfig{
					Namespace:   svc.Namespace,
					ServiceName: svc.Name,
					IP:          nip,
					Port:        p.NodePort,
					TargetPort:  p.TargetPort,
					Protocol:    string(p.Protocol),
					PortType:    cluster.PORT_TYPE_NODE_PORT,
				})
			}

		}
		if p.Port > 0 && len(svc.Spec.ExternalIPs) > 0 {
			//开放了EIP EPort
			for _, eip := range svc.Spec.ExternalIPs {
				cluster.GetClusterConfig().Add(&cluster.IpPortConfig{
					Namespace:   svc.Namespace,
					ServiceName: svc.Name,
					IP:          eip,
					Port:        p.Port,
					TargetPort:  p.TargetPort,
					Protocol:    string(p.Protocol),
					PortType:    cluster.PORT_TYPE_EIP_PORT,
				})
			}
		}
	}
}
