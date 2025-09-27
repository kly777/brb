package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"brb/internal/handler"
	"brb/internal/middleware"
	"brb/internal/repo"
	"brb/internal/router"
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

	userRepo, err := repo.NewUserRepo(a.DB)
	if err != nil {
		return fmt.Errorf("failed to create user repository: %w", err)
	}

	// 初始化services
	signService := service.NewSignService(signRepo)
	todoService := service.NewTodoService(todoRepo, taskRepo, eventRepo)
	taskService := service.NewTaskService(taskRepo, todoRepo)
	eventService := service.NewEventService(eventRepo, taskRepo)
	userService := service.NewUserService(userRepo)

	// 获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production" // 默认密钥，生产环境应该使用环境变量
		logger.Info.Println("使用默认JWT密钥，生产环境请设置JWT_SECRET环境变量")
	}

	// 初始化handlers
	signHandler := handler.NewSignHandler(signService)
	todoHandler := handler.NewTodoHandler(todoService)
	taskHandler := handler.NewTaskHandler(taskService)
	eventHandler := handler.NewEventHandler(eventService)
	userHandler := handler.NewUserHandler(userService, jwtSecret)

	// 创建路由注册器
	reg := router.NewStandardRouter(a.Mux)

	// 全局中间件
	reg.Use(
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware,
	)

	// API 版本分组
	v1 := reg.Group("/v1")

	// 注册公开路由（无需认证）
	signHandler.RegisterRoutes(v1)
	userHandler.RegisterRoutes(v1)

	// 为受保护的路由组添加认证中间件
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(jwtSecret))

	// 注册受保护的路由
	todoHandler.RegisterRoutes(protected)
	taskHandler.RegisterRoutes(protected)
	eventHandler.RegisterRoutes(protected)

	return nil
}

func (a *App) Run(addr string) error {

	// 创建自定义的 HTTP 服务器配置
	server := &http.Server{
		Addr:    addr,
		Handler: a.Mux,
		// 超时时间
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
