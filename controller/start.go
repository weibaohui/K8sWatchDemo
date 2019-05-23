package controller

import (
	"K8sWatchDemo/handler"
)

func Run() {

	var conf = &Config{
		Handlers: handler.Map,
		Resource: Resource{
			Deployment:            false,
			ReplicationController: false,
			ReplicaSet:            false,
			DaemonSet:             false,
			Services:              true,
			Pod:                   true,
			Job:                   false,
			PersistentVolume:      false,
			Namespace:             false,
			Secret:                false,
			ConfigMap:             false,
			Ingress:               false,
		},
		Namespace: "",
	}
	Start(conf)
}
