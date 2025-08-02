package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"brb/internal/sign"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func Inject() {
	// 初始化SQLite数据库
	db, err := sql.Open("sqlite3", "file:signs.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// 初始化sign存储
	signStore, err := sign.NewSignRepo(db)
	if err != nil {
		log.Fatalf("failed to create sign store: %v", err)
	}

	// 初始化sign处理器
	signHandler := sign.NewSignHandler(signStore)

	// 创建路由
	router := mux.NewRouter()
	signHandler.RegisterRoutes(router.PathPrefix("/api").Subrouter())

	// 启动HTTP服务器
	fmt.Println("服务器运行在 :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
