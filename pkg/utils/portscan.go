package utils

import (
	"K8sWatchDemo/pkg/cluster"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

func CheckPort(address string) (bool, error) {
	conn, e := net.DialTimeout("tcp4", address, time.Millisecond*500)
	if e != nil {
		logrus.Error(e.Error())
		return false, e
	} else {
		logrus.Info(address, "OPEN")
		conn.Close()
	}
	return true, nil
}

func checkUpdate(v *cluster.IpPortConfig) {
	//todo 增加一个list 更新方法，改为go 异步执行
	result, _ := CheckPort(fmt.Sprintf("%s:%d", v.IP, v.Port))
	if result {
		v.Linkable = true
	} else {
		v.Linkable = false
	}
	v.LastLinkTime = time.Now().Format("2006-01-02 15:04:05")
}

func AutoCheckPorts() {
	// tick := time.NewTicker(time.Minute * 2)
	// for {
	// 	select {
	// 	case <-tick.C:
	// 		for _, v := range cluster.GetClusterConfig().List {
	// 			go checkUpdate(v)
	// 		}
	// 	}
	// }

}
