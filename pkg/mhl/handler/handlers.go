package handler

import (
	"K8sWatchDemo/pkg/mhl/handler/deploy"
	"K8sWatchDemo/pkg/mhl/handler/pod"
	"K8sWatchDemo/pkg/watcher"
	"k8s.io/client-go/informers"
)

func Register(f informers.SharedInformerFactory) error {
	r := &watcher.Register{
		Handlers: []watcher.HandlersRegister{
			deploy.Register,
			pod.Register,
		},
	}
	return r.Register(f)
}
