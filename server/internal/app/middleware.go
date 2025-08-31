package app

import (
	"net/http"
	"net/url"
	"strings"

	"brb/pkg/logger"
)

// Middleware 中间件函数类型
type Middleware func(http.Handler) http.Handler

// setCORSHeaders 设置CORS头
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// handlePreflight 处理预检请求
func handlePreflight(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

// CreateProxyHandler 创建代理处理器
func CreateProxyHandler(mux *http.ServeMux, targetURL *url.URL) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		setCORSHeaders(w)

		// 处理预检请求
		if handlePreflight(w, r) {
			return
		}

		// 如果是API请求，由Mux处理

		logger.Info.Println("API请求:", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})
}

// CreateLoggingDivider 创建带日志分隔符的中间件
func CreateLoggingDivider(handler http.Handler) http.Handler {
	divider := strings.Repeat("-", 50)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println(divider)
		handler.ServeHTTP(w, r)
	})
}
