package main

import (
	v1 "k8s.io/api/core/v1"
	"strings"
)

func getPodName(podNameNs string) (namespace, podName string) {
	if strings.Contains(podNameNs, "/") {
		names := strings.SplitN(podNameNs, "/", 2)
		namespace = names[0]
		podName = names[1]
		return
	}
	return "", podNameNs
}

func getSvcName(podName string) string {
	uid := getCommonUID(podName)
	svcName := podSelectorName + "-svc-" + uid
	// fmt.Printf("SVC NAME:%s \n", svcName)
	return svcName
}

// 获取通用名称
func getCommonUID(podName string) string {
	after := strings.SplitN(podName, "-", 2)
	if len(after) == 2 {
		return after[1]
	}
	return ""
}

// 设置PodName为label
func addPodNameLabels(pod *v1.Pod) bool {

	oldLabels := pod.GetLabels()
	// 没有 PodName
	if oldLabels["podName"] == "" {
		labels := make(map[string]string)
		for e := range oldLabels {
			labels[e] = oldLabels[e]
		}
		labels["podName"] = pod.Name
		pod.SetLabels(labels)
		return true
	}
	return false
}
