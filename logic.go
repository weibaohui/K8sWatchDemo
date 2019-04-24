package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strings"
)

func checkPod(podSelectorName string, clientset *kubernetes.Clientset) {
	podList, e := clientset.CoreV1().Pods("default").List(metav1.ListOptions{
		LabelSelector: "app=" + podSelectorName,
	})
	if e != nil {
		panic(e.Error())
	}
	fmt.Printf("dubbo app pod 共有 %d 个\n", len(podList.Items))
	for i := range podList.Items {
		pod := podList.Items[i]

		after := strings.SplitN(pod.Name, "-", 2)
		if len(after) == 2 {
			serviceName := podSelectorName + "-svc-" + after[1]
			checkService(serviceName, clientset)
		}
	}
}

func checkService(svcName string, clientset *kubernetes.Clientset) {
	fmt.Println("svcName=", svcName)
	serviceList, e := clientset.CoreV1().Services("default").List(metav1.ListOptions{
		FieldSelector: "metadata.name=" + svcName,
	})
	if e != nil {
		fmt.Println("获取Service失败", e.Error())
	}
	if len(serviceList.Items) >= 0 {
		e := clientset.CoreV1().Services("default").Delete(svcName, &metav1.DeleteOptions{})
		if e == nil {
			fmt.Println("service 删除 成功", svcName)
		}
	}

	fmt.Println("准备创建 service ", svcName)

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": "dubbo-svc",
			},
			Ports: []v1.ServicePort{
				{Name: "web", Port: 8080, TargetPort: intstr.FromInt(80)},
				{Name: "test", Port: 8081, TargetPort: intstr.FromInt(81)},
			},
		},
	}

	create, e := clientset.CoreV1().Services("default").Create(svc)
	if e != nil {
		fmt.Println("创建 service 失败", e.Error())

	}
	fmt.Println("创建 service 成功", create.Name)
}
