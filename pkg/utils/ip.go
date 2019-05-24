package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/weibaohui/go-kit/httpkit"
	"time"
)

var ingressIPs []string
var nodeIPs []string

func IngressIPs() []string {
	return ingressIPs
}

func NodeIPs() []string {
	return nodeIPs
}

func SyncIPConfig() {
	str, err := httpkit.Get("http://127.0.0.1:9999/").String()
	if err != nil {
		logrus.Error(err.Error())
	}
	str = `["192.168.110.245","192.168.110.246"]`
	err = json.Unmarshal([]byte(str), &ingressIPs)
	err = json.Unmarshal([]byte(str), &nodeIPs)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func AutoSyncIPConfig() {
	tick := time.NewTicker(time.Minute * 30)
	for {
		select {
		case <-tick.C:
			SyncIPConfig()
		}
	}
}
