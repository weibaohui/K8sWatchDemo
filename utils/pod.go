package utils

import v1 "k8s.io/api/core/v1"

// 设置PodName为label
func AddPodNameLabels(pod *v1.Pod) bool {
	// 不存在 podName 更新
	if !IsPodLabelExists(pod, "podName") {
		oldLabels := pod.GetLabels()
		labels := make(map[string]string)
		for e := range oldLabels {
			labels[e] = oldLabels[e]
		}
		labels["podName"] = pod.Name
		pod.SetLabels(labels)
		return true
	}

	// 存在 podName 不更新
	return false
}

func IsPodLabelExists(pod *v1.Pod, k string) bool {
	labels := pod.GetLabels()
	return labels[k] != ""
}
