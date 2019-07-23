package pod

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type Starter struct {
}

func (s *Starter) Start(f informers.SharedInformerFactory, stop chan struct{}) error {
	go f.Core().V1().Pods().Informer().Run(stop)
	return nil
}
func (s *Starter) Sync(f informers.SharedInformerFactory, stop chan struct{}) error {
	hasSynced := f.Core().V1().Pods().Informer().HasSynced
	if synced := cache.WaitForCacheSync(stop, hasSynced); synced {
		logrus.Info("Pods 同步完成")
	}
	return nil
}
