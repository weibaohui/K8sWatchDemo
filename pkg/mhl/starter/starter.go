package starter

import (
	"context"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/informers"
)

type Starter interface {
	Sync(f informers.SharedInformerFactory, stop chan struct{}) error
	Start(f informers.SharedInformerFactory, stop chan struct{}) error
}

func StartThenSync(f informers.SharedInformerFactory, stop chan struct{}, starters ...Starter) error {
	if err := Start(f, stop, starters...); err != nil {
		return err
	}
	return Sync(f, stop, starters...)
}

func Sync(f informers.SharedInformerFactory, stop chan struct{}, starters ...Starter) error {
	eg, _ := errgroup.WithContext(context.TODO())
	for _, starter := range starters {
		func(starter Starter) {
			eg.Go(func() error {
				return starter.Sync(f, stop)
			})
		}(starter)
	}
	return eg.Wait()
}

func Start(f informers.SharedInformerFactory, stop chan struct{}, starters ...Starter) error {
	for _, starter := range starters {
		if err := starter.Start(f, stop); err != nil {
			return err
		}
	}
	return nil
}
