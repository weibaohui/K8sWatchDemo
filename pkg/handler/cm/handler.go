package cm

import (
	"K8sWatchDemo/pkg/cluster"
	"K8sWatchDemo/pkg/event"
	"K8sWatchDemo/pkg/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

type ConfigMapHandler struct {
	logger *logrus.Entry
}

func (h *ConfigMapHandler) Init() {
	h.logger = logrus.WithField("handler", "ConfigMapHandler")
}

//cm 监控的是整个CM的变更，不是某一个内部条目
func (h *ConfigMapHandler) ObjectCreated(obj interface{}) {
	cm := obj.(*v1.ConfigMap)
	if cm.Name != "tcp-services" && cm.Name != "udp-services" {
		return
	}
	if len(cm.Data) == 0 {
		return
	}

	protocol := "TCP"
	if cm.Name == "udp-services" {
		protocol = "UDP"
	}

	//先删除所有的
	cluster.GetClusterConfig().ClearIngressSvc(protocol)
	addIngressSvc(cm.Data, protocol)

	h.logger.Infoln("add", cm.Name)
}
func (h *ConfigMapHandler) ObjectDeleted(event event.InformerEvent) {
	//要删除所有的
	if event.Name != "tcp-services" && event.Name != "udp-services" {
		return
	}

	protocol := "TCP"
	if event.Name == "udp-services" {
		protocol = "UDP"
	}

	cluster.GetClusterConfig().ClearIngressSvc(protocol)

	h.logger.Infoln("delete", event)
}

func (h *ConfigMapHandler) ObjectUpdated(oldObj interface{}, event event.InformerEvent) {
	cm := oldObj.(*v1.ConfigMap)
	if cm.Name != "tcp-services" && cm.Name != "udp-services" {
		return
	}

	protocol := "TCP"
	if cm.Name == "udp-services" {
		protocol = "UDP"
	}

	cluster.GetClusterConfig().ClearIngressSvc(protocol)

	addIngressSvc(cm.Data, protocol)

	h.logger.Infoln("update", cm.Name)
}

func addIngressSvc(data map[string]string, protocol string) {
	// 8886 default/svc-2:8886
	for _, v := range data {
		ns, svcName, port, err := getNsNamePort(v)
		if err != nil {
			return
		}

		for _, ip := range utils.IngressIPs() {

			cluster.GetClusterConfig().Add(&cluster.IpPortConfig{
				Namespace:      ns,
				IngressSvcName: svcName,
				IP:             ip,
				Port:           port,
				PortType:       cluster.PORT_TYPE_INGRESS_NGINX_PORT,
				Protocol:       protocol,
			})
		}

	}
}

func getNsNamePort(item string) (ns, svcName string, port int32, err error) {
	//default/svc-2:8886
	names := strings.SplitN(item, "/", 2)
	if len(names) != 2 {
		err = errors.New("格式错误" + item)
		return
	}
	if len(names) == 2 {
		ns = names[0]
	}

	nameports := strings.SplitN(names[1], ":", 2)
	if len(nameports) != 2 {
		err = errors.New("格式错误" + item)
		return
	}
	svcName = nameports[0]
	i, err := strconv.Atoi(nameports[1])
	if err != nil {
		return
	}

	port = int32(i)
	return
}
