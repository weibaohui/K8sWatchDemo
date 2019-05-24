package pkg

import (
	"K8sWatchDemo/pkg/config"
	"K8sWatchDemo/pkg/controller"
	"K8sWatchDemo/pkg/handler"
)

func Start() {

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
	controller.Start(conf)
}
