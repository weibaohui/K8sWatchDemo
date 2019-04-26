package watcher

import (
	"fmt"
	"strings"
	"sync"
)

// 监控目标
type target struct {
	Namespace     string
	PodNamePrefix string
}

func isTargetByPodName(namespace, podName string) bool {
	for _, v := range targets {
		if v.Namespace == namespace &&
			strings.HasPrefix(podName, v.PodNamePrefix+"-") {
			return true
		}
	}
	return false
}

func isTarget(namespace, podName string, f func(string, string) bool) bool {
	result := f(namespace, podName)
	fmt.Printf("目标 %s %s 监听: %t \n", namespace, podName, result)
	return result
}

var targets map[string]*target
var once sync.Once

func init() {
	once.Do(func() {
		fmt.Println("初始化监控对象容器")
		targets = make(map[string]*target)
	})
}

func AddTarget(namespace, podNamePrefix string) {
	key := fmt.Sprintf("%s-%s", namespace, podNamePrefix)
	targets[key] = &target{
		Namespace:     namespace,
		PodNamePrefix: podNamePrefix,
	}
	fmt.Printf("%s 增加成功，开始监控\n\n", key)
}

func DeleteTarget(namespace, podNamePrefix string) {
	// todo 应该调用处理逻辑，删除所有相关的svc，ep
	key := fmt.Sprintf("%s-%s", namespace, podNamePrefix)
	delete(targets, key)
	fmt.Printf("%s 删除,执行清理工作\n\n", key)
}
