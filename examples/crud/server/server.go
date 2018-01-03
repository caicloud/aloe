package server

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

type Product struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func (p Product) get(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	resp.WriteHeaderAndEntity(200, Product{
		Id:    id,
		Title: "test",
	})
}

func (p Product) post(req *restful.Request, resp *restful.Response) {
	updatedProduct := new(Product)
	err := req.ReadEntity(updatedProduct)
	if err != nil { // bad request
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	resp.WriteHeaderAndEntity(201, updatedProduct)
}

func (p Product) delete(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	log.Println("getting product with id:" + id)
	resp.WriteHeader(204)
}

// Register register route
func (p Product) Register() {
	ws := new(restful.WebService)
	ws.Path("/products")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{id}").To(p.get).
		Doc("get the product by its id").
		Param(ws.PathParameter("id", "identifier of the product").DataType("string")))

	ws.Route(ws.POST("").To(p.post).
		Doc("create a product").
		Param(ws.BodyParameter("Product", "a Product").DataType("main.Product")))

	ws.Route(ws.PUT("").To(p.post).
		Doc("update a product").
		Param(ws.BodyParameter("Product", "a Product").DataType("main.Product")))

	ws.Route(ws.DELETE("/{id}").To(p.delete).
		Doc("delete a product").
		Param(ws.BodyParameter("Product", "identifier of the product").DataType("string")))

	restful.Add(ws)
}
