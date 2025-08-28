package main

import (
	"log"

	"brb/internal/app"

)

func main() {
	// 创建应用程序实例
	app, err := app.NewApp("file:brb.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// 启动HTTP服务器
	log.Fatal(app.Run(":8080"))
}
