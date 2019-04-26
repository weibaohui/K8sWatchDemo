package pkg

import (
	"K8sWatchDemo/utils"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Action struct {
	PodNameNs  string // default/dubbo-xxxxx-xx
	ActionName string
}

// 删除逻辑
func (h *Helper) DeleteProcess(podNameNs string) {
	ns, svcName, _ := utils.GetNsSvcPodName(podNameNs)
	// 删除Service
	h.deleteSvc(ns, svcName)
}

// 更新POD的处理逻辑
// 每次POD的状态更新都会触发
// 如
func (h *Helper) UpdateProcess(podNameNs string) {
	ns, svcName, podName := utils.GetNsSvcPodName(podNameNs)

	if pod, e := h.GetPod(ns, podName); e == nil {
		h.addPodNameToLabelIfAbsent(pod)

		// 查看更新的状态
		for _, v := range pod.Status.ContainerStatuses {
			if v.State.Waiting != nil {
				println(pod.Name, v.Name, v.State.Waiting.Reason)
			}
		}

		// 根据POD 状态 创建或者删除SVC
		if h.isPodReady(pod) {
			// POD READY 且没有创建对应的SVC,应创建
			if !h.IsServiceExists(ns, svcName) {
				fmt.Printf("POD %s 已经ready,缺SVC,为其创建SVC \n", pod.Name)
				h.createSvc(ns, podName)
			}
		} else {
			// POD NOT READY,如果有SVC,应删除SVC
			if h.IsServiceExists(ns, svcName) {
				fmt.Printf("POD %s 没有ready,但已经为其创建了SVC,删除 \n", pod.Name)
				h.deleteSvc(ns, svcName)
			}

		}
	}

}

// 新增POD的处理逻辑，第一次启动初始化时也会进入次程序
func (h *Helper) AddProcess(pod *v1.Pod) {
	// 如果程序初始化运行，会收到所有已经存在的POD，应该先检查POD状态

	// 检查podName 是否设置了，更新podName
	h.addPodNameToLabelIfAbsent(pod)

	ns, podName := pod.Namespace, pod.Name
	svcName := utils.GetSvcName(podName)

	// not ready 但是已经有svc，删除，如果pod变回ready 会触发更新事件
	if pod, e := h.GetPod(ns, podName); e == nil {
		if !h.isPodReady(pod) && h.IsServiceExists(ns, svcName) {
			fmt.Println("POD NOT READY,删除SVC", podName)
			h.deleteSvc(ns, svcName)
		}
	}
}

// 检查POD是否READY，
func (h *Helper) isPodReady(pod *v1.Pod) bool {
	for e := range pod.Status.ContainerStatuses {
		status := pod.Status.ContainerStatuses[e]
		// 如果pod 没有准备好，应该删除svc
		if status.Ready == false {
			return false
		}
	}
	return true
}

func (h *Helper) createSvc(ns, podName string) {
	svcName := utils.GetSvcName(podName)
	config, e := getNetConfig(podName)
	if e != nil {
		fmt.Println("创建 SVC 失败", svcName)
		fmt.Println(e.Error())
		return
	}

	fmt.Println(config)

	svc := &v1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			// ClusterIP: v1.ClusterIPNone,
			ExternalIPs: config.ExternalIPs,
			Selector: map[string]string{
				"podName": podName,
			},
			Ports: config.ServicePorts,
		},
	}

	create, e := h.Services(ns).Create(svc)
	if e != nil {
		fmt.Println("创建 service 失败", e.Error())

	} else {
		fmt.Println("创建 service 成功", create.Name)
	}

}

func (h *Helper) deleteSvc(ns, svcName string) {
	fmt.Println("删除 SVC", svcName)
	if e := h.Services(ns).Delete(svcName, &metaV1.DeleteOptions{}); e != nil {
		fmt.Printf("删除 SVC %s 失败:%s\n", svcName, e.Error())
	}
}
