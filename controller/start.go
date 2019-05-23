package controller

import (
	"K8sWatchDemo/config"
	"K8sWatchDemo/handler"
)

func Run() {

	var conf = &config.Config{
		Handlers: handler.Map,
		Resource: config.Resource{
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
