package utils

import (
	v1 "k8s.io/api/core/v1"
	"strings"
)

func GetPodName(podNameNs string) (namespace, podName string) {
	if strings.Contains(podNameNs, "/") {
		names := strings.SplitN(podNameNs, "/", 2)
		namespace = names[0]
		podName = names[1]
		return
	}
	return "", podNameNs
}

func GetSvcName(podName string) string {
	prefix, uid := GetCommonUID(podName)
	svcName := prefix + "-svc-" + uid
	// fmt.Printf("SVC NAME:%s \n", svcName)
	return svcName
}

// 获取通用名称
// todo 能否加入其他格式的处理？
func GetCommonUID(podName string) (prefix, uid string) {
	names := strings.SplitN(podName, "-", 2)
	if len(names) == 2 {
		prefix = names[0]
		uid = names[1]
		return
	}
	return "", ""
}

// 设置PodName为label
func AddPodNameLabels(pod *v1.Pod) bool {

	oldLabels := pod.GetLabels()
	// 没有 PodNameNs
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
