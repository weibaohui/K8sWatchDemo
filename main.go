package main

import (
	"K8sWatchDemo/pkg"
	"K8sWatchDemo/pkg/utils"
	"K8sWatchDemo/pkg/webservice"
	"fmt"

	"sync"
	"time"
)

var nodeports = sync.Map{}
var eipports = sync.Map{}

func main() {
	utils.SyncIPConfig()

	go pkg.Start()

	go webservice.Start()

	go utils.AutoCheckPorts()
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
