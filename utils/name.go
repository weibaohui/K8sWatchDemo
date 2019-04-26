package utils

import (
	"strings"
)

func GetNsSvcPodName(podNameNs string) (ns, svcName, podName string) {
	ns, podName = GetPodName(podNameNs)
	svcName = GetSvcName(podName)
	return ns, svcName, podName

}

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
