package webservice

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)

func Start() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/index/{id}").
		To(index).
		Produces(restful.MIME_JSON).
		Doc("主页"))
	restful.Add(ws)
	log.Fatal(http.ListenAndServe(":9999", nil))

}

// GET /users/1
func index(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	type name struct {
		Name string `json:"name"`
		Age  int    `json:"agex"`
		Log  int    `json:"log,omitempty"`
	}
	response.WriteEntity(name{
		Name: id,
		Age:  56,
		Log:  0,
	})
}
