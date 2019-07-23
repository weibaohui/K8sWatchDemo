package deployment

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type Starter struct {
}

func (s *Starter) Start(f informers.SharedInformerFactory, stop chan struct{}) error {
	go f.Apps().V1().Deployments().Informer().Run(stop)
	return nil
}
func (s *Starter) Sync(f informers.SharedInformerFactory, stop chan struct{}) error {
	hasSynced := f.Apps().V1().Deployments().Informer().HasSynced
	if synced := cache.WaitForCacheSync(stop, hasSynced); synced {
		logrus.Info("Deployment 同步完成")
	}
	return nil
}
