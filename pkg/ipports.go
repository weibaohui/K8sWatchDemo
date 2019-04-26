package pkg

import (
	"fmt"
	"github.com/weibaohui/go-kit/httpkit"
)

type ipPort struct {
}

// 获取IP 端口 对应关系
// todo 改为实际的
func getIpPort() (*ipPort, error) {
	s, e := httpkit.Get("http://www.baidu.com").String()
	if e != nil {
		return nil, e
	}
	fmt.Println("模拟在线获取", len(s))
	return &ipPort{}, nil
}
