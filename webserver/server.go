package webserver

import (
	"K8sWatchDemo/watcher"
	"fmt"
	"net/http"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("ns")
	pod := r.URL.Query().Get("pod")
	if len(ns) > 0 && len(pod) > 0 {
		watcher.AddTarget(ns, pod)
		fmt.Fprintln(w, "添加成功 "+ns+pod)
	} else {
		fmt.Fprintln(w, "请输入 ns pod")

	}

}

func DelHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("ns")
	pod := r.URL.Query().Get("pod")
	if len(ns) > 0 && len(pod) > 0 {
		watcher.DeleteTarget(ns, pod)
		fmt.Fprintln(w, "删除成功 "+ns+pod)
	} else {
		fmt.Fprintln(w, "请输入 ns pod")

	}

}
func Start() {
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/del", DelHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)

}
