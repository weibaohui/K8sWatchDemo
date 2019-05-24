package main

import (
	"K8sWatchDemo/pkg"
	"fmt"

	"sync"
	"time"
)

var nodeports = sync.Map{}
var eipports = sync.Map{}

func main() {
	go pkg.Start()

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
