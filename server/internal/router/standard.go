package router

import (
	"brb/internal/middleware"
	"net/http"
	"strings"
)

// standardRouter 适配 net/http.ServeMux
type standardRouter struct {
	prefix string
	mux    *http.ServeMux
	mws    []middleware.Middleware
}

func NewStandardRouter(mux *http.ServeMux) Router {
	return &standardRouter{
		mux: mux,
	}
}

func (r *standardRouter) GET(path string, handler http.HandlerFunc, mws ...middleware.Middleware) {
	r.handle("GET", path, handler, mws...)
}

func (r *standardRouter) POST(path string, handler http.HandlerFunc, mws ...middleware.Middleware) {
	r.handle("POST", path, handler, mws...)
}

func (r *standardRouter) PUT(path string, handler http.HandlerFunc, mws ...middleware.Middleware) {
	r.handle("PUT", path, handler, mws...)
}

func (r *standardRouter) DELETE(path string, handler http.HandlerFunc, mws ...middleware.Middleware) {
	r.handle("DELETE", path, handler, mws...)
}

func (r *standardRouter) PATCH(path string, handler http.HandlerFunc, mws ...middleware.Middleware) {
	r.handle("PATCH", path, handler, mws...)
}

func (r *standardRouter) handle(method, path string, handler http.HandlerFunc, mws ...middleware.Middleware) {
	fullPath := r.prefix + path

	mws = append(r.mws, mws...)
	// 应用中间件
	h := http.Handler(handler)
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	// 注册路由
	r.mux.Handle(strings.ToUpper(method)+" "+fullPath, h)
}

func (r *standardRouter) Group(prefix string) Router {
	return &standardRouter{
		prefix: r.prefix + prefix,
		mux:    r.mux,
		mws:    r.mws,
	}
}

func (r *standardRouter) Use(mws ...middleware.Middleware) {
	r.mws = append(r.mws, mws...)
}
