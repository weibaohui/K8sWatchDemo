package main

import (
	"K8sWatchDemo/pkg/mhl"
)

func main() {

	go mhl.Start()
	select {}

}
