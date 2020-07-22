package webservice

import (
	"K8sWatchDemo/pkg/cluster"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)

func Start() {
	container := restful.NewContainer()
	ws := new(restful.WebService)
	ws.Route(ws.GET("/ports").
		To(ports).
		Produces(restful.MIME_JSON))
	container.Add(ws)
	log.Fatal(http.ListenAndServe(":9999", container))
}

// GET /ports
func ports(request *restful.Request, response *restful.Response) {
	configs := cluster.GetClusterConfig().List
	response.WriteEntity(configs)
}
