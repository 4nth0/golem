package router

import (
	"net/http"

	"github.com/4nth0/golem/router/tree"
)

type Router struct {
	Tree *tree.TreeNode
}

type RouteHandler func(w http.ResponseWriter, r *http.Request, params map[string]string)

type RouteHandlers map[string]interface{}

func NewRouter() *Router {
	return &Router{
		Tree: tree.NewTree(),
	}
}

func (r *Router) Get(path string, handler tree.Handler) {
	r.Add("GET", path, handler)
}

func (r *Router) Post(path string, handler tree.Handler) {
	r.Add("POST", path, handler)
}

func (r *Router) Put(path string, handler tree.Handler) {
	r.Add("PUT", path, handler)
}

func (r *Router) Delete(path string, handler tree.Handler) {
	r.Add("DELETE", path, handler)
}

func (r *Router) Add(method string, path string, handler tree.Handler) {
	r.Tree.AddNode(path, method, handler)
}

func (r *Router) GetHandler(path, method string) (tree.Handler, map[string]string, error) {
	handler, params, err := r.Tree.GetNode(path, method)

	return handler, params, err
}
