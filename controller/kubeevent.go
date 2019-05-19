package controller

// k8s 资源变更事件
type KubeEvent struct {
	Kind      string
	Name      string
	Namespace string
}
