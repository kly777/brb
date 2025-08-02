package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"brb/internal/handler"
	"brb/internal/repo"
	"brb/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// App 表示应用程序
type App struct {
	DB     *sql.DB
	Router *mux.Router
}

// NewApp 创建并初始化应用程序
func NewApp(dbPath string) (*App, error) {
	app := &App{
		Router: mux.NewRouter(),
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

	// 初始化services
	signService := service.NewSignService(signRepo)

	// 初始化handlers
	signHandler := handler.NewSignHandler(signService)

	// 注册路由
	apiRouter := a.Router.PathPrefix("/api").Subrouter()
	signHandler.RegisterRoutes(apiRouter)

	return nil
}

// Run 启动应用程序
func (a *App) Run(addr string) error {
	fmt.Printf("服务器运行在 %s\n", addr)
	return http.ListenAndServe(addr, a.Router)
}

// Close 关闭应用程序资源
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
