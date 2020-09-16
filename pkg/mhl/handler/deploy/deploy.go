package deploy

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	listerv1 "k8s.io/client-go/listers/apps/v1"
	"strings"
)

type DeploymentModule struct {
	cache listerv1.DeploymentLister
}

func register(f informers.SharedInformerFactory) {
	m := &DeploymentModule{}
	m.cache = f.Apps().V1().Deployments().Lister()
	informer := f.Apps().V1().Deployments().Informer()
	informer.AddEventHandler(m)
}

func (m *DeploymentModule) OnAdd(obj interface{}) {
	deployment, err := m.cache.Deployments("default").Get("nginx-deployment")
	if err != nil {
		logrus.Info(err.Error())
		if strings.Contains(err.Error(), "not found") {
			logrus.Info("cache 没有获取到", "nginx-deployment")
		}
	} else {
		logrus.Infof("从 cache 获取到 %s,有%v个副本", deployment.Name, *deployment.Spec.Replicas)
	}
	logrus.Infof("deploymentEventHandler OnAdd ,%v ", obj.(*v1.Deployment).Name)
}

func (m *DeploymentModule) OnUpdate(oldObj, newObj interface{}) {
	old := oldObj.(*v1.Deployment)
	newobj := newObj.(*v1.Deployment)
	if old.ObjectMeta.GetResourceVersion() == newobj.ObjectMeta.GetResourceVersion() {
		logrus.Info("same deploy %s,%s \n", old.Name, old.ObjectMeta.GetResourceVersion())
		return
	}
	logrus.Infof("deploymentEventHandler OnUpdate new ,%v,%s ", newobj.Name, newobj.ObjectMeta.GetResourceVersion())
}

func (m *DeploymentModule) OnDelete(obj interface{}) {
	logrus.Infof("deploymentEventHandler OnDelete ,%v ", obj)
}
