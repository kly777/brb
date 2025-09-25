package middleware

import (
	"net/http"
	"time"

	"brb/pkg/logger"
)

type Middleware func(http.Handler) http.Handler

// LoggingMiddleware 记录请求日志的中间件
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// RecoveryMiddleware 恢复panic的中间件
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error.Printf("Panic: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

