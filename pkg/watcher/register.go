package watcher

type Register struct {
	Handlers []HandlersRegister
}

func (r *Register) Register(w *Watcher, stop chan struct{}) error {
	for _, reg := range r.Handlers {
		 reg(w, stop)
	}
	return nil
}

type HandlersRegister func(w *Watcher, stop chan struct{})
