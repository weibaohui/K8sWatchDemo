package deploy

import (
	"K8sWatchDemo/pkg/watcher"
	"fmt"
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
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "not found") {
			deployment, err = m.cache.Deployments("default").Get("ratings-v1")
			if err != nil {
				fmt.Println("第二次又失败了", err.Error())
			}
		}
	} else {
		fmt.Println("cache 获取到的", deployment.Name)
	}
	fmt.Printf("deploymentEventHandler OnAdd ,%v \n\n", obj.(*v1.Deployment).Name)
}

func (m *DeploymentModule) OnUpdate(oldObj, newObj interface{}) {
	old := oldObj.(*v1.Deployment)
	newobj := newObj.(*v1.Deployment)
	if old.ObjectMeta.GetResourceVersion() == newobj.ObjectMeta.GetResourceVersion() {
		fmt.Printf("same deploy %s,%s \n", old.Name, old.ObjectMeta.GetResourceVersion())
		return
	}
	fmt.Printf("deploymentEventHandler OnUpdate new ,%v,%s \n\n", newobj.Name, newobj.ObjectMeta.GetResourceVersion())
}

func (m *DeploymentModule) OnDelete(obj interface{}) {
	fmt.Printf("deploymentEventHandler OnDelete ,%v \n\n", obj)
}
