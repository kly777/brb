package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"brb/internal/dto"
	"brb/internal/entity"
	"brb/internal/router"

	"github.com/golang-jwt/jwt/v5"
)

// userHandler 处理用户相关的HTTP请求
type userHandler struct {
	userService UserService
	jwtSecret   string
}

type UserService interface {
	Register(username, password string) (*entity.User, error)
	Login(username, password string) (*entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	GetAllUsers() ([]*entity.User, error)
	UpdateUser(id uint, username, password string, role entity.Role) (*entity.User, error)
	DeleteUser(id uint) error
	ChangePassword(id uint, oldPassword, newPassword string) error
	PromoteToAdmin(id uint) error
	DemoteToUser(id uint) error
}

// NewUserHandler 创建新的UserHandler
func NewUserHandler(userService UserService, jwtSecret string) *userHandler {
	return &userHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

// Register 用户注册
func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证请求数据
	if req.Username == "" || req.Password == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 生成JWT token
	token, err := h.generateJWT(user)
	if err != nil {
		http.Error(w, "生成token失败", http.StatusInternalServerError)
		return
	}

	response := dto.LoginResponse{
		User:  dto.FromUserEntity(user),
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login 用户登录
func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证请求数据
	if req.Username == "" || req.Password == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 生成JWT token
	token, err := h.generateJWT(user)
	if err != nil {
		http.Error(w, "生成token失败", http.StatusInternalServerError)
		return
	}

	response := dto.LoginResponse{
		User:  dto.FromUserEntity(user),
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser 获取当前用户信息
func (h *userHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "未授权访问", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.FromUserEntity(user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllUsers 获取所有用户（仅管理员）
func (h *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// 检查用户角色
	userRole, ok := r.Context().Value("userRole").(entity.Role)
	if !ok || userRole != entity.RoleAdmin {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}

	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := dto.FromUserEntities(users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// UpdateUser 更新用户信息
func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 检查权限：用户只能更新自己的信息，管理员可以更新任何用户
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "未授权访问", http.StatusUnauthorized)
		return
	}

	userRole, _ := r.Context().Value("userRole").(entity.Role)
	
	// 获取要更新的用户ID
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "无效的用户ID格式", http.StatusBadRequest)
		return
	}

	// 普通用户只能更新自己的信息
	if userRole != entity.RoleAdmin && userID != uint(targetID) {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}

	var req dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	// 普通用户不能修改角色
	if userRole != entity.RoleAdmin && req.Role != "" {
		http.Error(w, "普通用户不能修改角色", http.StatusForbidden)
		return
	}

	user, err := h.userService.UpdateUser(uint(targetID), req.Username, req.Password, entity.Role(req.Role))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := dto.FromUserEntity(user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteUser 删除用户
func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 只有管理员可以删除用户
	userRole, ok := r.Context().Value("userRole").(entity.Role)
	if !ok || userRole != entity.RoleAdmin {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "无效的用户ID格式", http.StatusBadRequest)
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ChangePassword 修改密码
func (h *userHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// 用户只能修改自己的密码
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "未授权访问", http.StatusUnauthorized)
		return
	}

	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "旧密码和新密码不能为空", http.StatusBadRequest)
		return
	}

	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PromoteToAdmin 提升用户为管理员
func (h *userHandler) PromoteToAdmin(w http.ResponseWriter, r *http.Request) {
	// 只有管理员可以提升用户权限
	userRole, ok := r.Context().Value("userRole").(entity.Role)
	if !ok || userRole != entity.RoleAdmin {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "无效的用户ID格式", http.StatusBadRequest)
		return
	}

	if err := h.userService.PromoteToAdmin(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DemoteToUser 降级用户为普通用户
func (h *userHandler) DemoteToUser(w http.ResponseWriter, r *http.Request) {
	// 只有管理员可以降级用户权限
	userRole, ok := r.Context().Value("userRole").(entity.Role)
	if !ok || userRole != entity.RoleAdmin {
		http.Error(w, "权限不足", http.StatusForbidden)
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "无效的用户ID格式", http.StatusBadRequest)
		return
	}

	if err := h.userService.DemoteToUser(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// generateJWT 生成JWT token
func (h *userHandler) generateJWT(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"userID": user.ID,
		"role":   string(user.Role),
		"exp":    jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// RegisterRoutes 注册用户相关路由
func (h *userHandler) RegisterRoutes(r router.Router) {
	// 公开路由（无需认证）
	public := r.Group("/api/auth")
	public.POST("/register", h.Register)
	public.POST("/login", h.Login)

	// 受保护的路由（需要认证）
	protected := r.Group("/api/users")
	
	// 个人操作路由
	protected.GET("/me", h.GetCurrentUser)
	protected.PUT("/password", h.ChangePassword)
	protected.PUT("/:id", h.UpdateUser)

	// 管理员操作路由
	protected.GET("", h.GetAllUsers)
	protected.DELETE("/:id", h.DeleteUser)
	protected.POST("/:id/promote", h.PromoteToAdmin)
	protected.POST("/:id/demote", h.DemoteToUser)
}