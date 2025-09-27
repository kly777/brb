package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"brb/internal/entity"
	"brb/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
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

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtSecret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 获取Authorization头
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "缺少认证令牌", http.StatusUnauthorized)
				return
			}

			// 检查Bearer token格式
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "认证令牌格式错误", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// 解析和验证JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// 验证签名方法
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// 提取用户信息
				userID, ok := claims["userID"].(float64)
				if !ok {
					http.Error(w, "令牌中缺少用户ID", http.StatusUnauthorized)
					return
				}

				role, ok := claims["role"].(string)
				if !ok {
					http.Error(w, "令牌中缺少用户角色", http.StatusUnauthorized)
					return
				}

				// 将用户信息添加到请求上下文
				ctx := context.WithValue(r.Context(), "userID", uint(userID))
				ctx = context.WithValue(ctx, "userRole", entity.Role(role))
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
			}
		})
	}
}

// RequireAuth 要求认证的中间件（包装AuthMiddleware，简化使用）
func RequireAuth(jwtSecret string) Middleware {
	return AuthMiddleware(jwtSecret)
}

// RequireRole 要求特定角色的中间件
func RequireRole(requiredRole entity.Role) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从上下文中获取用户角色
			userRole, ok := r.Context().Value("userRole").(entity.Role)
			if !ok {
				http.Error(w, "未授权访问", http.StatusUnauthorized)
				return
			}

			// 检查用户角色是否满足要求
			if userRole != requiredRole {
				http.Error(w, "权限不足", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin 要求管理员角色的中间件
func RequireAdmin() Middleware {
	return RequireRole(entity.RoleAdmin)
}

// OptionalAuth 可选认证中间件（不强制要求认证，但如果有token会解析）
func OptionalAuth(jwtSecret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString := parts[1]
					
					token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
						}
						return []byte(jwtSecret), nil
					})

					if err == nil && token.Valid {
						if claims, ok := token.Claims.(jwt.MapClaims); ok {
							if userID, ok := claims["userID"].(float64); ok {
								ctx := context.WithValue(r.Context(), "userID", uint(userID))
								if role, ok := claims["role"].(string); ok {
									ctx = context.WithValue(ctx, "userRole", entity.Role(role))
								}
								r = r.WithContext(ctx)
							}
						}
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
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

