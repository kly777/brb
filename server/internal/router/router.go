package router

import (
	"brb/internal/middleware"
	"net/http"
)

// Router 是路由注册器接口
type Router interface {
    // 基本HTTP方法
    GET(path string, handler http.HandlerFunc)
    POST(path string, handler http.HandlerFunc)
    PUT(path string, handler http.HandlerFunc)
    DELETE(path string, handler http.HandlerFunc)
    PATCH(path string, handler http.HandlerFunc)
    
    // 路由分组
    Group(prefix string) Router
    
    // 中间件支持
    Use(middlewares ...middleware.Middleware)
}