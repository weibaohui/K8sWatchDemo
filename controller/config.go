package controller

import "K8sWatchDemo/handler"

// Resource contains resource configuration
type Resource struct {
	Deployment            bool `json:"deploy"`
	ReplicationController bool `json:"rc"`
	ReplicaSet            bool `json:"rs"`
	DaemonSet             bool `json:"ds"`
	Services              bool `json:"svc"`
	Pod                   bool `json:"po"`
	Job                   bool `json:"job"`
	PersistentVolume      bool `json:"pv"`
	Namespace             bool `json:"ns"`
	Secret                bool `json:"secret"`
	ConfigMap             bool `json:"cm"`
	Ingress               bool `json:"ing"`
}

type Config struct {
	Handlers  map[string]handler.Handler `json:"handlers"`
	Resource  Resource                   `json:"resource"`
	Namespace string                     `json:"namespace,omitempty"`
}
