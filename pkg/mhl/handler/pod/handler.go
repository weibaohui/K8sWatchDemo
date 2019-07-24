package pod

import (
	"k8s.io/client-go/informers"
)

func Register(f informers.SharedInformerFactory) {
	register(f)
}
