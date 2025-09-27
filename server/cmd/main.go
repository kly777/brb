package main

import (
	"brb/internal/app"
	"brb/pkg/logger"
)

func main() {
	logger.Info.Println("Starting server...")
	// 创建应用程序实例
	app, err := app.NewApp("file:brb.db?cache=shared&mode=rwc")
	if err != nil {
		logger.Error.Fatalf("Failed to initialize application: %v", err)
	}

	// 启动HTTP服务器
	err = app.Run(":5050")
	if err != nil {
		logger.Error.Fatalf("Failed to start server: %v", err)
	}
}
