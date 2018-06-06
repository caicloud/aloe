package server

import (
	"fmt"
	"log"

	"github.com/caicloud/aloe/runtime"
	"github.com/emicklei/go-restful"
)

// ProductServer defines a product server
type ProductServer struct {
	ps map[string]Product
}

// NewProductServer returns a product server
func NewProductServer() *ProductServer {
	return &ProductServer{
		ps: map[string]Product{},
	}
}

// Name implements cleaner.Cleaner
func (s *ProductServer) Name() string {
	return "product"
}

// Clean implements cleaner.Cleaner
func (s *ProductServer) Clean(template *runtime.RoundTripTemplate, args map[string]string) error {
	if args == nil {
		return fmt.Errorf("missing args")
	}
	v, ok := args["product"]
	if !ok {
		return fmt.Errorf("variable product is not registered")
	}
	delete(s.ps, v)
	return nil
}

// Product defines product example
type Product struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (s *ProductServer) get(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	p, ok := s.ps[id]
	if !ok {
		resp.WriteHeader(404)
		return
	}
	if err := resp.WriteEntity(p); err != nil {
		log.Printf("response error: %v", err)
	}
}

func (s *ProductServer) create(req *restful.Request, resp *restful.Response) {
	p := Product{}
	err := req.ReadEntity(&p)
	if err != nil {
		resp.WriteHeader(400)
		return
	}
	if _, ok := s.ps[p.ID]; ok {
		resp.WriteHeader(409)
		return
	}
	s.ps[p.ID] = p
	if err := resp.WriteHeaderAndEntity(201, p); err != nil {
		log.Printf("response error: %v", err)
	}
}

func (s *ProductServer) update(req *restful.Request, resp *restful.Response) {
	p := Product{}
	err := req.ReadEntity(&p)
	if err != nil {
		resp.WriteHeader(400)
		return
	}
	id := req.PathParameter("id")
	s.ps[id] = p
	if err := resp.WriteHeaderAndEntity(200, p); err != nil {
		log.Printf("response error: %v", err)
	}
}

func (s *ProductServer) remove(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	if _, ok := s.ps[id]; ok {
		delete(s.ps, id)
	}
	resp.WriteHeader(204)
}

// Register register route
func (s *ProductServer) Register() {
	ws := new(restful.WebService)
	ws.Path("/products")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{id}").To(s.get).
		Doc("get the product by its id").
		Param(ws.PathParameter("id", "identifier of the product").DataType("string")))

	ws.Route(ws.POST("").To(s.create).
		Doc("create a product").
		Param(ws.BodyParameter("Product", "a Product").DataType("main.Product")))

	ws.Route(ws.PUT("/{id}").To(s.update).
		Doc("update a product").
		Param(ws.BodyParameter("Product", "a Product").DataType("main.Product")))

	ws.Route(ws.DELETE("/{id}").To(s.remove).
		Doc("delete a product").
		Param(ws.BodyParameter("Product", "identifier of the product").DataType("string")))

	restful.Add(ws)
}
