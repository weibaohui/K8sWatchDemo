package controller

import (
	"log"
)

func Run() {
	var eventHandler = new(DefaultHandler)

	var conf = &Config{
		Resource: Resource{
			Deployment:            true,
			ReplicationController: true,
			ReplicaSet:            true,
			DaemonSet:             true,
			Services:              true,
			Pod:                   true,
			Job:                   true,
			PersistentVolume:      true,
			Namespace:             true,
			Secret:                true,
			ConfigMap:             true,
			Ingress:               true,
		},
		Namespace: "",
	}
	if err := eventHandler.Init(conf); err != nil {
		log.Fatal(err)
	}
	Start(conf, eventHandler)
}
