package pod

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	corev1 "k8s.io/client-go/listers/core/v1"
)

type PodModule struct {
	cache corev1.PodLister
}

func register(f informers.SharedInformerFactory) {
	m := &PodModule{}
	m.cache = f.Core().V1().Pods().Lister()
	informer := f.Core().V1().Pods().Informer()
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
