package deploy

import (
	"K8sWatchDemo/pkg/watcher"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	listerv1 "k8s.io/client-go/listers/apps/v1"
	"strings"
)

type DeploymentModule struct {
	watcher *watcher.Watcher
	cache   listerv1.DeploymentLister
}

func register(w *watcher.Watcher, stop chan struct{}) {
 	m := &DeploymentModule{}
	m.watcher = w
	m.cache = m.watcher.Deployments.Lister()
	informer := m.watcher.Deployments.Informer()
	informer.AddEventHandler(m)
}

func (m *DeploymentModule) OnAdd(obj interface{}) {
	deployment, err := m.cache.Deployments("default").Get("ratings-v1")
	if err != nil {
		logrus.Info(err.Error())
		if strings.Contains(err.Error(), "not found") {
			deployment, err = m.cache.Deployments("default").Get("ratings-v1")
			if err != nil {
				logrus.Info("第二次又失败了", err.Error())
			}
		}
	} else {
		logrus.Info("cache 获取到的", deployment.Name)
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
