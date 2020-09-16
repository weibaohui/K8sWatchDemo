package deploy

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	listerv1 "k8s.io/client-go/listers/apps/v1"
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
	// deployment, err := m.cache.Deployments("default").Get("busybox")
	// if err != nil {
	// 	logrus.Info(err.Error())
	// 	// if strings.Contains(err.Error(), "not found") {
	// 	// 	deployment, err = m.cache.Deployments("default").Get("ratings-v1")
	// 	// 	if err != nil {
	// 	// 		logrus.Info("第二次又失败了", err.Error())
	// 	// 	}
	// 	// }
	// } else {
	// 	logrus.Info("cache 获取到的", deployment.Name)
	// }
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
