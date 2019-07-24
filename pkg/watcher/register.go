package watcher

import "k8s.io/client-go/informers"

type Register struct {
	Handlers []HandlersRegister
}

func (r *Register) Register(f informers.SharedInformerFactory) error {
	for _, reg := range r.Handlers {
		reg(f)
	}
	return nil
}

type HandlersRegister func(f informers.SharedInformerFactory)
