package app

import (
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"

	"brb/pkg/logger"
	"maps"
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

// isAPIRequest 检查是否为API请求
func isAPIRequest(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/api/")
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
		if isAPIRequest(r) {
			logger.Info.Println("API请求:", r.URL.Path)
			mux.ServeHTTP(w, r)
			return
		}

		// 否则，代理到前端开发服务器
		logger.Info.Printf("请求前端: %s\n", r.URL.Path)
		fetchAndServeHTML(w, r, targetURL)
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

// fetchAndServeHTML 获取HTML内容并处理
func fetchAndServeHTML(w http.ResponseWriter, r *http.Request, targetURL *url.URL) {
	// 构造目标请求
	proxyReq, err := http.NewRequest(r.Method, targetURL.ResolveReference(r.URL).String(), r.Body)
	if err != nil {
		logger.Error.Printf("Error constructing request: %v", err)
		http.Error(w, "Failed to construct request", http.StatusInternalServerError)
		return
	}

	// 复制请求头
	proxyReq.Header = make(http.Header)
	maps.Copy(proxyReq.Header, r.Header)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		logger.Error.Printf("Error fetching HTML: %v", err)
		http.Error(w, "Failed to fetch HTML", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	logger.Info.Printf("Fetched content: %s", string(body[:20]))
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// 检查 Content-Type 是否为 HTML
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		logger.Info.Printf("Non-HTML content type: %s", contentType)
		// 非 HTML 内容，直接返回原始响应
		// 复制原始响应头
		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		logger.Info.Println("Returned non-HTML content directly")
		return
	}

	// 解析HTML模板
	tmpl, err := template.New("frontend").Parse(string(body))
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	// 准备模板数据（示例数据）
	data := struct {
		Title string
	}{
		Title: "Title From Backend",
	}

	// 执行模板并写入响应
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}
