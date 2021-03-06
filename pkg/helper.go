package pkg

import (
	"K8sWatchDemo/utils"
	"flag"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typeV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sync"
)

var cli *kubernetes.Clientset
var once = sync.Once{}

type Helper struct {
	cli *kubernetes.Clientset
}

func NewHelper() *Helper {
	once.Do(func() {
		cli = getClient()
	})

	return &Helper{cli: cli}
}

func (h *Helper) RESTClient() rest.Interface {
	return h.cli.CoreV1().RESTClient()
}
func (h *Helper) Pods(ns string) typeV1.PodInterface {
	return h.cli.CoreV1().Pods(ns)
}
func (h *Helper) Services(ns string) typeV1.ServiceInterface {
	return h.cli.CoreV1().Services(ns)
}

func (h *Helper) GetPod(ns, podName string) (*coreV1.Pod, error) {
	return h.Pods(ns).Get(podName, metaV1.GetOptions{})
}

func (h *Helper) GetService(ns, svcName string) (*coreV1.Service, error) {
	return h.Services(ns).Get(svcName, metaV1.GetOptions{})
}
func (h *Helper) IsServiceExists(ns, svcName string) bool {
	_, e := h.Services(ns).Get(svcName, metaV1.GetOptions{})
	if e == nil {
		return true
	}
	return false
}

func (h *Helper) addPodNameToLabelIfAbsent(pod *coreV1.Pod) {
	// 检查podName 是否设置了，更新podName
	if utils.AddPodNameLabels(pod) {
		pod, e := h.Pods(pod.Namespace).Update(pod)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		fmt.Println("增加 PodNameNs 到 metadata.labels", pod.Name)
	}
}
func getClient() *kubernetes.Clientset {
	var kubeConfig *string
	var inCluster *bool
	if home := homeDir(); home != "" {
		s := filepath.Join(home, ".kube", "config")
		kubeConfig = flag.String("kubeConfig", s, "kubeconfig存放位置")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "kubeconfig存放位置")
	}
	inCluster = flag.Bool("in", false, "是否在集群内")
	flag.Parse()
	var config *rest.Config
	var err error
	if *inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	}
	if err != nil {
		panic(err.Error())
	}
	cli, e := kubernetes.NewForConfig(config)
	if e != nil {
		panic(e.Error())
	}
	return cli

}

func homeDir() string {
	if s := os.Getenv("HOME"); s != "" {
		return s
	}
	return os.Getenv("USERPROFILE")
}
