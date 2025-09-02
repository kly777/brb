package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
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
	todoService := service.NewTodoService(todoRepo, taskRepo, eventRepo)
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

	// 使用中间件处理器
	proxyHandler := CreateProxyHandler(a.Mux, target)
	loggingHandler := CreateLoggingDivider(proxyHandler)

	// 创建自定义的 HTTP 服务器配置
	server := &http.Server{
		Addr:    addr,
		Handler: loggingHandler,
		// 设置合理的超时时间
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info.Printf("服务器运行在 %s\n", addr)
	defer func() {
		a.Close()
		logger.Info.Println("服务器已关闭")
	}()

	// 启动服务器
	return server.ListenAndServe()
}

// Close 关闭应用程序资源
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
