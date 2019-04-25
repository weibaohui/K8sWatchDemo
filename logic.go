package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
)

type Action struct {
	PodName    string
	ActionName string
}

const podSelectorName = "dubbo"

func getPodName(podNameNs string) (namespace, podName string) {
	if strings.Contains(podNameNs, "/") {
		names := strings.SplitN(podNameNs, "/", 2)
		namespace = names[0]
		podName = names[1]
		return
	}
	return "", podNameNs
}

func isTarget(podName string) bool {
	return strings.HasPrefix(podName, podSelectorName+"-")
}
func (h *helper) deletePodProcess(podName string) {
	// 删除Service
	e := h.deleteSvc(podName)
	if e != nil {
		fmt.Println(e.Error())
	}

}

func (h *helper) updatePodProcess(podName string) {
	h.updatePodSelector(podName)
}
func (h *helper) addPodProcess(podName string) {
	h.updatePodSelector(podName)

	if !h.isServiceExists(podName) {
		h.createSvc(podName)
	}
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

func (h *helper) updatePodSelector(podName string) {
	pod, e := h.Pods("default").Get(podName, metav1.GetOptions{})
	if e != nil {
		fmt.Println(" 无此POD ", e.Error())
		return
	}

	// 增加了PodName Label，再更新
	if addPodNameLabels(pod) {
		_, e = h.Pods("default").Update(pod)
		if e != nil {
			fmt.Println(e.Error())
		}
		fmt.Println("增加 PodName 到 metadata.labels", podName)
	}

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

func (h *helper) isServiceExists(podName string) bool {
	svcName := getSvcName(podName)
	serviceList, e := h.Services("default").List(metav1.ListOptions{
		FieldSelector: "metadata.name=" + svcName,
	})

	if e != nil {
		fmt.Println(e.Error())
		return false
	}

	if len(serviceList.Items) > 0 {
		return true
	}

	return false
}

func (h *helper) createSvc(podName string) {
	svcName := getSvcName(podName)
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			// ClusterIP: v1.ClusterIPNone,
			Selector: map[string]string{
				"podName": podName,
			},
			Ports: []v1.ServicePort{
				{Name: "web", Port: 8080, TargetPort: intstr.FromInt(80)},
				{Name: "test", Port: 8081, TargetPort: intstr.FromInt(81)},
			},
		},
	}
	create, e := h.Services("default").Create(svc)
	if e != nil {
		fmt.Println("创建 service 失败", e.Error())

	} else {
		fmt.Println("创建 service 成功", create.Name)
	}

}

func (h *helper) deleteSvc(podName string) error {
	svcName := getSvcName(podName)
	fmt.Println("删除 SVC", svcName)
	return h.Services("default").Delete(svcName, &metav1.DeleteOptions{})

}
