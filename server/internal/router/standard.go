package router

import (
	"brb/internal/middleware"
	"net/http"
	"strings"
)

// standardRouter 适配 net/http.ServeMux
type standardRouter struct {
	prefix      string
	mux         *http.ServeMux
	middlewares []middleware.Middleware
}

func NewStandardRouter(mux *http.ServeMux) Router {
	return &standardRouter{
		mux: mux,
	}
}

func (r *standardRouter) GET(path string, handler http.HandlerFunc) {
	r.handle("GET", path, handler)
}

func (r *standardRouter) POST(path string, handler http.HandlerFunc) {
	r.handle("POST", path, handler)
}

func (r *standardRouter) PUT(path string, handler http.HandlerFunc) {
	r.handle("PUT", path, handler)
}

func (r *standardRouter) DELETE(path string, handler http.HandlerFunc) {
	r.handle("DELETE", path, handler)
}

func (r *standardRouter) PATCH(path string, handler http.HandlerFunc) {
	r.handle("PATCH", path, handler)
}

func (r *standardRouter) handle(method, path string, handler http.HandlerFunc) {
	fullPath := r.prefix + path

	// 应用中间件
	h := http.Handler(handler)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	// 注册路由
	r.mux.Handle(strings.ToUpper(method)+" "+fullPath, h)
}

func (r *standardRouter) Group(prefix string) Router {
	return &standardRouter{
		prefix:      r.prefix + prefix,
		mux:         r.mux,
		middlewares: r.middlewares,
	}
}

func (r *standardRouter) Use(middlewares ...middleware.Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}
