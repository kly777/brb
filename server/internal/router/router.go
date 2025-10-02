package router

import (
	"brb/internal/middleware"
	"net/http"
)

// Router 是路由注册器接口
type Router interface {
	// 基本HTTP方法
	GET(path string, handler http.HandlerFunc, middlewares ...middleware.Middleware)
	POST(path string, handler http.HandlerFunc, middlewares ...middleware.Middleware)
	PUT(path string, handler http.HandlerFunc, middlewares ...middleware.Middleware)
	DELETE(path string, handler http.HandlerFunc, middlewares ...middleware.Middleware)
	PATCH(path string, handler http.HandlerFunc, middlewares ...middleware.Middleware)

	// 路由分组
	Group(prefix string) Router

	// 中间件支持
	Use(middlewares ...middleware.Middleware)
}
