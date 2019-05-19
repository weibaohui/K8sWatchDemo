package main

import (
	"K8sWatchDemo/controller"
	"fmt"

	"sync"
	"time"
)

var nodeports = sync.Map{}
var eipports = sync.Map{}

func main() {
	go controller.Run()

	// go ApiWatchStart()
	// go printUsedPorts(&nodeports)
	// go printUsedPorts(&eipports)
	select {}

}

func printUsedPorts(ports *sync.Map) {
	for {
		time.Sleep(time.Second * 5)
		ports.Range(func(key, value interface{}) bool {
			fmt.Printf("%v  ", key)
			return true
		})
		fmt.Println()
	}

}
