package watcher

import (
	"fmt"
	"strings"
)

func isTargetByPodName(namespace, podName string) bool {
	return strings.HasPrefix(podName, PodSelectorName+"-")
}

func isTarget(namespace, podName string, f func(string, string) bool) bool {
	result := f(namespace, podName)
	fmt.Println("是否监听目标", namespace, podName, result)
	return result
}
