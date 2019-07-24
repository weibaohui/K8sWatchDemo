package mhl

import (
	"K8sWatchDemo/pkg/mhl/handler"
	"K8sWatchDemo/pkg/mhl/starter"
	"K8sWatchDemo/pkg/mhl/starter/deployment"
	"K8sWatchDemo/pkg/mhl/starter/pod"
	"K8sWatchDemo/pkg/utils"
	"K8sWatchDemo/pkg/watcher"
	"k8s.io/client-go/informers"
	"time"
)

func Start() {
	stop := make(chan struct{})
	var kubeClient = utils.NewHelper().GetKubeClient()

	factory := informers.NewSharedInformerFactory(kubeClient, time.Hour*2)

	w := watcher.Watcher{
		Factory:      factory,
		DaemonSets:   factory.Apps().V1().DaemonSets(),
		Deployments:  factory.Apps().V1().Deployments(),
		ReplicaSets:  factory.Apps().V1().ReplicaSets(),
		StatefulSets: factory.Apps().V1().StatefulSets(),
		ConfigMaps:   factory.Core().V1().ConfigMaps(),
		Endpoints:    factory.Core().V1().Endpoints(),
		Namespaces:   factory.Core().V1().Namespaces(),
		Pods:         factory.Core().V1().Pods(),
		Services:     factory.Core().V1().Services(),
		Ingresses:    factory.Extensions().V1beta1().Ingresses(),
		Nodes:        factory.Core().V1().Nodes(),
	}

	//
	w.Factory.Start(stop)

	w.Factory.WaitForCacheSync(stop)

	starter.StartThenSync(factory, stop, &deployment.Starter{}, &pod.Starter{})

	handler.Register(factory)

	select {}
}
