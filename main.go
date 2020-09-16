package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/weibaohui/go-kit/strkit"
	v1 "k8s.io/api/apps/v1"
	api_v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"K8sWatchDemo/pkg/mhl"
	"K8sWatchDemo/pkg/utils"
)

var nodeports = sync.Map{}
var eipports = sync.Map{}

func main() {
	// utils.SyncIPConfig()

	// go pkg.Start()

	// go webservice.Start()

	// go utils.AutoCheckPorts()
	// go ApiWatchStart()
	// go printUsedPorts(&nodeports)
	// go printUsedPorts(&eipports)

	go mhl.Start()
	client := utils.NewHelper().GetKubeClient()

	for i := 0; i < 4000; i++ {
		_, err := client.AppsV1().Deployments("default").Create(newDeployment())
		if err != nil {
			logrus.Error(err.Error())
		}
		_, err =client.CoreV1().Pods("default").Create(newPod())
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	select {}

}

func newPod() *api_v1.Pod {
	return &api_v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test--pod---" + strkit.RandomString(6),
			Namespace: "default",
		},
		Spec: api_v1.PodSpec{
			Containers: []api_v1.Container{
				{
					Name:  "nginx",
					Image: "nginx:alpine",
				},
			},
		},
	}
}

func newDeployment() *v1.Deployment {
	var replicas int32
	replicas = 2
	d := &v1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test--deploy--" + strkit.RandomString(6),
			Namespace: "default",
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,

			Template: api_v1.PodTemplateSpec{
				Spec: api_v1.PodSpec{
					Containers: []api_v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:alpine",
						},
					},
				},
			},
		},
	}
	return d
}

func printUsedPorts(ports *sync.Map) {
	for {
		time.Sleep(time.Second * 5)
		ports.Range(func(key, value interface{}) bool {
			fmt.Printf("%v  ", key)
			return true
		})
		fmt.Println()
	}

}
