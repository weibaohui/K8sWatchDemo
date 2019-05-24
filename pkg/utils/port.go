package utils

import (
	"K8sWatchDemo/pkg/cluster"
	"fmt"
	"net"
	"strings"
	"time"
)

func Check(network, address string) (bool, error) {
	conn, e := net.DialTimeout(network, address, time.Second*3)
	if e != nil {
		fmt.Println(e.Error())
		return false, e
	}

	defer conn.Close()
	return true, nil
}

func AutoCheck() {
	tick := time.NewTicker(time.Minute * 1)

	for {
		select {
		case <-tick.C:
			fmt.Println("2分钟一次，到点啦")
			for _, v := range cluster.GetClusterConfig().List {
				result, e := Check(strings.ToLower(v.Protocol), fmt.Sprintf("%s:%d", v.IP, v.Port))
				if result &&  e != nil {
					v.Linkable = true
				} else {
					v.Linkable = false
					fmt.Println(e.Error())
				}
				v.LastLinkTime = time.Now().Format("2006-01-02 15:04:05")
			}
		}
	}

}
