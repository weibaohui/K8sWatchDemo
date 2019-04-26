package pkg

import (
	"K8sWatchDemo/utils"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Action struct {
	PodNameNs  string // default/dubbo-xxxxx-xx
	ActionName string
}

// 删除逻辑
func (h *Helper) DeleteProcess(podNameNs string) {
	ns, podName := utils.GetPodName(podNameNs)

	// 删除Service
	e := h.deleteSvc(ns, podName)
	if e != nil {
		fmt.Println(e.Error())
	}

}

// 更新POD的处理逻辑
func (h *Helper) UpdateProcess(podNameNs string) {
	ns, podName := utils.GetPodName(podNameNs)
	h.updatePodSelector(ns, podName)
}

// 新增POD的处理逻辑
func (h *Helper) AddProcess(podNameNs string) {
	ns, podName := utils.GetPodName(podNameNs)
	h.updatePodSelector(ns, podName)

	if svc, e := h.isServiceExists(ns, podName); e == nil && svc == nil {
		h.createSvc(ns, podName)
	}
}

func (h *Helper) updatePodSelector(ns, podName string) {
	if pod, e := h.isPodExists(ns, podName); e == nil {
		// 增加了PodName Label，再更新
		if utils.AddPodNameLabels(pod) {
			_, e = h.Pods(ns).Update(pod)
			if e != nil {
				fmt.Println(e.Error())
			}
			fmt.Println("增加 PodNameNs 到 metadata.labels", podName)
		}
	}

}

func (h *Helper) isServiceExists(ns, podName string) (*v1.Service, error) {
	svcName := utils.GetSvcName(podName)
	list, e := h.Services(ns).List(metav1.ListOptions{
		FieldSelector: "metadata.name=" + svcName,
		Limit:         1,
	})
	if e != nil {
		return nil, e
	}

	if len(list.Items) == 0 {
		return nil, e
	}

	return &list.Items[0], nil
}

func (h *Helper) isPodExists(ns, podName string) (*v1.Pod, error) {
	return h.Pods(ns).Get(podName, metav1.GetOptions{})
}

func (h *Helper) createSvc(ns, podName string) {
	svcName := utils.GetSvcName(podName)
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
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

	create, e := h.Services(ns).Create(svc)
	if e != nil {
		fmt.Println("创建 service 失败", e.Error())

	} else {
		fmt.Println("创建 service 成功", create.Name)
	}

}

func (h *Helper) deleteSvc(ns, podName string) error {
	svcName := utils.GetSvcName(podName)
	fmt.Println("删除 SVC", svcName)
	return h.Services(ns).Delete(svcName, &metav1.DeleteOptions{})
}
