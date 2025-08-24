package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"brb/internal/handler"
	"brb/internal/repo"
	"brb/internal/service"

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

	// 初始化依赖
	if err := app.initDependencies(); err != nil {
		return nil, err
	}

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
	// 添加CORS中间件
	corsMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用原始的多路复用器
		a.Mux.ServeHTTP(w, r)
	})

	fmt.Printf("服务器运行在 %s\n", addr)
	defer func() { log.Println("服务器已关闭") }()
	return http.ListenAndServe(addr, corsMux)
}

// Close 关闭应用程序资源
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}

