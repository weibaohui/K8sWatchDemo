package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strings"
)

const (
	ADD    string = "ADD"
	DELETE string = "DELETE"
	UPDATE string = "UPDATE"
	SYNC   string = "SYNC"
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
func deletePodProcess(cli *kubernetes.Clientset, podName string) {
	// 删除Service
	e := deleteSvc(cli, getSvcName(podSelectorName, podName))
	if e != nil {
		fmt.Println(e.Error())
	}

}

func addPodProcess(cli *kubernetes.Clientset, podName string) {
	svcName := getSvcName(podSelectorName, podName)
	checkServiceExists(cli, svcName)
	createSvc(cli, svcName)
}

func getSvcName(podSelectorName, podName string) string {
	uid := getCommonUID(podName)
	svcName := podSelectorName + "-svc-" + uid
	fmt.Printf("SVC NAME:%s \n", svcName)
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

func checkPod(clientset *kubernetes.Clientset, podSelectorName string) {
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
			checkServiceExists(clientset, serviceName)
		}
	}
}

func checkServiceExists(cli *kubernetes.Clientset, svcName string) bool {
	serviceList, e := cli.CoreV1().Services("default").List(metav1.ListOptions{
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

func createSvc(cli *kubernetes.Clientset, svcName string) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: v1.ServiceSpec{
			Type:      v1.ServiceTypeClusterIP,
			ClusterIP: v1.ClusterIPNone,
			Selector: map[string]string{
				"app": podSelectorName + "-svc",
			},
			Ports: []v1.ServicePort{
				{Name: "web", Port: 8080, TargetPort: intstr.FromInt(80)},
				{Name: "test", Port: 8081, TargetPort: intstr.FromInt(81)},
			},
		},
	}
	create, e := cli.CoreV1().Services("default").Create(svc)
	if e != nil {
		fmt.Println("创建 service 失败", e.Error())

	} else {
		fmt.Println("创建 service 成功", create.Name)
	}

}

func deleteSvc(cli *kubernetes.Clientset, svcName string) error {
	fmt.Println("删除 SVC", svcName)
	return cli.CoreV1().Services("default").Delete(svcName, &metav1.DeleteOptions{})

}