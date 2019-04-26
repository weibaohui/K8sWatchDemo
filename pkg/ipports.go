package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/weibaohui/go-kit/strkit"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

type NetConfig struct {
	ExternalIPs  []string
	ServicePorts []v1.ServicePort
}

// 获取IP 端口 对应关系
// todo 改为实际的
func getNetConfig(podName string) (*NetConfig, error) {

	fmt.Printf("获取 %fakeStr 的 IPPORT配置信息 \n", podName)

	config := NetConfig{}
	fakeStr := `{
  "ExternalIPs": [
    "192.168.1.1",
    "192.168.1.2"
  ],
  "ServicePorts": [
    {
      "name": "web2",
      "protocol": "TCP",
      "port": 8082,
      "targetPort": 80,
      "nodePort": 31384
    },
    {
      "name": "test3",
      "protocol": "UDP",
      "port": 8083,
      "targetPort": 81,
      "nodePort": 31357
    }
  ]
}`
	e := json.Unmarshal([]byte(fakeStr), &config)

	if e != nil {
		fmt.Println(e.Error())
		return nil, e
	}

	// todo 加入 ExternalIPs ServicePorts 校验
	config = NetConfig{
		ExternalIPs: []string{"192.168.1.1", "192.168.1.2"},
		ServicePorts: []v1.ServicePort{
			{Name: "web2", Port: 8082, TargetPort: intstr.FromInt(80), NodePort: fakeNodePort(), Protocol: "TCP"},
			{Name: "test3", Port: 8083, TargetPort: intstr.FromInt(81), NodePort: fakeNodePort(), Protocol: "UDP"},
		},
	}

	return &config, nil

}

// 30000-32767
func fakeNodePort() int32 {
	i, _ := strconv.Atoi(strkit.RandomNumber(3))
	i2 := i + 31000
	return int32(i2)
}
