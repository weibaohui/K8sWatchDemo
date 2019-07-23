package pod

import (
	"K8sWatchDemo/pkg/watcher"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/client-go/listers/core/v1"
 )

type PodModule struct {
	watcher *watcher.Watcher
	cache   corev1.PodLister
}

func register(w *watcher.Watcher, stop chan struct{}) {
	m := &PodModule{}
	m.watcher = w
	m.cache = m.watcher.Pods.Lister()
	informer := m.watcher.Pods.Informer()
	informer.AddEventHandler(m)
}

func (m *PodModule) OnAdd(obj interface{}) {
	 
	logrus.Infof("podEventHandler OnAdd ,%v ", obj.(*v1.Pod).Name)
}

func (m *PodModule) OnUpdate(oldObj, newObj interface{}) {
	old := oldObj.(*v1.Pod)
	newobj := newObj.(*v1.Pod)
	if old.ObjectMeta.GetResourceVersion() == newobj.ObjectMeta.GetResourceVersion() {
		logrus.Info("same deploy %s,%s \n", old.Name, old.ObjectMeta.GetResourceVersion())
		return
	}
	logrus.Infof("podEventHandler OnUpdate new ,%v,%s ", newobj.Name, newobj.ObjectMeta.GetResourceVersion())
}

func (m *PodModule) OnDelete(obj interface{}) {
	logrus.Infof("podEventHandler OnDelete ,%v ", obj)
}
