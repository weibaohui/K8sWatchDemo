package handler

import (
	"K8sWatchDemo/pkg/mhl/handler/deploy"
	"K8sWatchDemo/pkg/mhl/handler/pod"
	"K8sWatchDemo/pkg/watcher"
)

func Register(w *watcher.Watcher, stop chan struct{}) error {
	r := &watcher.Register{
		Handlers: []watcher.HandlersRegister{
			deploy.Register,
			pod.Register,
		},
	}
	return r.Register(w, stop)
}
