package cluster

import (
	"k8s.io/apimachinery/pkg/util/intstr"
	"sync"
)

type clusterConfig struct {
	w    sync.RWMutex
	List []*IpPortConfig
}

const (
	PORT_TYPE_NODE_PORT          = "NodePort"
	PORT_TYPE_EIP_PORT           = "EIPPort"
	PORT_TYPE_INGRESS_NGINX_PORT = "IngressNginxPort"
)

type IpPortConfig struct {
	Namespace      string             `json:"namespace"`
	ServiceName    string             `json:"service_name"` // 一版SVC
	IngressSvcName string             `json:"ingress_name"` // ingress nginx L4 使用的SVC
	PodName        string             `json:"pod_name"`
	IP             string             `json:"ip"`
	Port           int32              `json:"port"`
	PortType       string             `json:"port_type"`
	TargetPort     intstr.IntOrString `json:"target_port"`
	Linkable       bool               `json:"linkable"`
}

var o sync.Once
var cnc *clusterConfig

func init() {
	o.Do(func() {
		cnc = &clusterConfig{
			List: make([]*IpPortConfig, 0),
		}
	})

}

func GetClusterConfig() *clusterConfig {
	return cnc
}

func (c *clusterConfig) Add(config *IpPortConfig) {
	c.w.Lock()
	defer c.w.Unlock()
	c.List = append(c.List, config)
}

// 删除 普通 SVC
func (c *clusterConfig) DeleteSvc(ns string, svcName string) {
	c.w.Lock()
	defer c.w.Unlock()
	if len(c.List) == 0 {
		//没有数据
		return
	}
	for k := 0; k < len(c.List); k++ {
		v := c.List[k]
		if v.Namespace == ns && v.ServiceName == svcName {
			//前面的不动，隔一个，再拼上后面的,k需要减1，因为后面的元素index,往前移动了一个
			c.List = append(c.List[:k], c.List[k+1:]...)
			k--
		}
	}
}

// 删除 Ingress-nginx L4 中使用的 svc
func (c *clusterConfig) ClearIngressSvc() {
	c.w.Lock()
	defer c.w.Unlock()
	if len(c.List) == 0 {
		//没有数据
		return
	}
	for k := 0; k < len(c.List); k++ {
		v := c.List[k]
		if v.PortType == PORT_TYPE_INGRESS_NGINX_PORT {
			//前面的不动，隔一个，再拼上后面的,k需要减1，因为后面的元素index,往前移动了一个
			c.List = append(c.List[:k], c.List[k+1:]...)
			k--
		}
	}
}
