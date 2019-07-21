package deploy

import "K8sWatchDemo/pkg/watcher"

func Register(w *watcher.Watcher, stop chan struct{}) {
	register(w, stop)
}
