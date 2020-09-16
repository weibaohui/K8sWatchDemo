package pod

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	corev1 "k8s.io/client-go/listers/core/v1"

	"K8sWatchDemo/pkg/cluster"
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
	pod := obj.(*v1.Pod)
	logrus.Infof("podEventHandler OnAdd ,%v ", pod.Name)
	cnc := cluster.GetClusterConfig()
	cnc.Add(pod)
}

func (m *PodModule) OnUpdate(oldObj, newObj interface{}) {
	oldPod := oldObj.(*v1.Pod)
	newPod := newObj.(*v1.Pod)
	cnc := cluster.GetClusterConfig()
	cnc.Update(newPod)
	if oldPod.ObjectMeta.GetResourceVersion() == newPod.ObjectMeta.GetResourceVersion() {
		logrus.Info("same deploy %s,%s \n", oldPod.Name, oldPod.ObjectMeta.GetResourceVersion())
		return
	}
	logrus.Infof("podEventHandler OnUpdate new ,%v,%s ", newPod.Name, newPod.ObjectMeta.GetResourceVersion())
}

func (m *PodModule) OnDelete(obj interface{}) {
	logrus.Infof("podEventHandler OnDelete ,%v ", obj)
}
