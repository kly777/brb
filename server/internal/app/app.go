package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"
	"time"

	"brb/internal/handler"
	"brb/internal/repo"
	"brb/internal/service"
	"brb/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

// App 表示应用程序
type App struct {
	DB  *sql.DB
	Mux *http.ServeMux
}

// NewApp 创建并初始化应用程序
func NewApp(dbPath string) (*App, error) {
	app := &App{
		Mux: http.NewServeMux(),
	}

	// 初始化数据库连接
	var err error
	app.DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	logger.Info.Println("数据库连接已打开")

	// 初始化依赖
	if err := app.initDependencies(); err != nil {
		return nil, err
	}
	logger.Info.Println("依赖注入完成")

	return app, nil
}

// initDependencies 初始化应用程序依赖
func (a *App) initDependencies() error {
	// 初始化repositories
	signRepo, err := repo.NewSignRepo(a.DB)
	if err != nil {
		return fmt.Errorf("failed to create sign repository: %w", err)
	}

	todoRepo, err := repo.NewTodoRepo(a.DB)
	if err != nil {
		return fmt.Errorf("failed to create todo repository: %w", err)
	}

	taskRepo, err := repo.NewTaskRepo(a.DB)
	if err != nil {
		return fmt.Errorf("failed to create task repository: %w", err)
	}

	eventRepo, err := repo.NewEventRepo(a.DB)
	if err != nil {
		return fmt.Errorf("failed to create event repository: %w", err)
	}

	// 初始化services
	signService := service.NewSignService(signRepo)
	todoService := service.NewTodoService(todoRepo, taskRepo)
	taskService := service.NewTaskService(taskRepo, todoRepo)
	eventService := service.NewEventService(eventRepo, taskRepo)

	// 初始化handlers
	signHandler := handler.NewSignHandler(signService)
	todoHandler := handler.NewTodoHandler(todoService)
	taskHandler := handler.NewTaskHandler(taskService)
	eventHandler := handler.NewEventHandler(eventService)

	// 注册路由 - 使用net/http的新特性，方法匹配和路径模式
	signHandler.RegisterRoutes(a.Mux)
	todoHandler.RegisterRoutes(a.Mux)
	taskHandler.RegisterRoutes(a.Mux)
	eventHandler.RegisterRoutes(a.Mux)

	return nil
}

// Run 启动应用程序
func (a *App) Run(addr string) error {
	// 创建反向代理到前端开发服务器
	target, _ := url.Parse("http://localhost:5173")

	// 添加CORS中间件和代理逻辑
	proxyMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 如果是API请求，由Mux处理
		if strings.HasPrefix(r.URL.Path, "/api/") {
			a.Mux.ServeHTTP(w, r)
			return
		}

		// 否则，代理到前端开发服务器
		logger.Info.Printf("请求前端: %s\n", r.URL.Path)
		a.fetchAndServeHTML(w, r, target)
	})

	logger.Info.Printf("服务器运行在 %s\n", addr)
	defer func() {
		logger.Info.Println("服务器已关闭")
	}()
	time.Sleep(3 * time.Second)
	return http.ListenAndServe(addr, proxyMux)
}

// Close 关闭应用程序资源
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}

// fetchAndServeHTML 获取HTML内容并处理
func (a *App) fetchAndServeHTML(w http.ResponseWriter, r *http.Request, targetURL *url.URL) {
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
	logger.Info.Printf("Fetched content: %s", string(body[:100]))
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

		w.Header().Set("Content-Type", contentType)
		logger.Info.Println("contentType", w.Header().Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
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
		Title: "GO Title",
	}

	// 执行模板并写入响应
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}
