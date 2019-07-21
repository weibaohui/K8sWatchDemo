package watcher

import (
	"k8s.io/client-go/informers"
	v1 "k8s.io/client-go/informers/apps/v1"
	v12 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/informers/extensions/v1beta1"
)

type Watcher struct {
	Factory      informers.SharedInformerFactory
	DaemonSets   v1.DaemonSetInformer
	Deployments  v1.DeploymentInformer
	ReplicaSets  v1.ReplicaSetInformer
	StatefulSets v1.StatefulSetInformer
	ConfigMaps   v12.ConfigMapInformer
	Endpoints    v12.EndpointsInformer
	Namespaces   v12.NamespaceInformer
	Pods         v12.PodInformer
	Services     v12.ServiceInformer
	Ingresses    v1beta1.IngressInformer
	Nodes        v12.NodeInformer
}
