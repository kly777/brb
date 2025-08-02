package main

import (
	"log"

	"brb/internal/app"
)

func main() {
	// 创建应用程序实例
	application, err := app.NewApp("file:signs.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Close()

	// 启动HTTP服务器
	log.Fatal(application.Run(":8080"))
}
