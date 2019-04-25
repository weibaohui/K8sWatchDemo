package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Action struct {
	PodName    string
	ActionName string
}

func (h *Helper) deletePodProcess(podName string) {
	// 删除Service
	e := h.deleteSvc(podName)
	if e != nil {
		fmt.Println(e.Error())
	}

}

func (h *Helper) updatePodProcess(podName string) {
	h.updatePodSelector(podName)
}
func (h *Helper) addPodProcess(podName string) {
	h.updatePodSelector(podName)

	if !h.isServiceExists(podName) {
		h.createSvc(podName)
	}
}

func (h *Helper) updatePodSelector(podName string) {
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

func (h *Helper) isServiceExists(podName string) bool {
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

func (h *Helper) createSvc(podName string) {
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

func (h *Helper) deleteSvc(podName string) error {
	svcName := getSvcName(podName)
	fmt.Println("删除 SVC", svcName)
	return h.Services("default").Delete(svcName, &metav1.DeleteOptions{})

}
