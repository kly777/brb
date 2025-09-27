package dto

import (
	"brb/internal/entity"
	"time"
)

// UserRegisterRequest 用户注册请求DTO
type UserRegisterRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	// 注册时不能指定角色，默认为user
}

// UserLoginRequest 用户登录请求DTO
type UserLoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// UserUpdateRequest 用户更新请求DTO
type UserUpdateRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
}

// UserResponse 用户响应DTO
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginResponse 登录响应DTO
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// ToEntity 将UserRegisterRequest转换为entity.User
func (req *UserRegisterRequest) ToEntity() *entity.User {
	user := &entity.User{
		Username:  req.Username,
		Password:  req.Password, // 注意：在实际使用中应该加密密码
		Role:      entity.RoleUser, // 注册时固定为user角色
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user
}

// ToEntity 将UserUpdateRequest转换为entity.User
func (req *UserUpdateRequest) ToEntity(id uint) *entity.User {
	user := &entity.User{
		ID:        id,
		Username:  req.Username,
		Password:  req.Password, // 注意：在实际使用中应该加密密码
		UpdatedAt: time.Now(),
	}

	if req.Role != "" {
		user.Role = entity.Role(req.Role)
	}

	return user
}

// FromUserEntity 将entity.User转换为UserResponse
func FromUserEntity(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// FromUserEntities 将entity.User切片转换为UserResponse切片
func FromUserEntities(users []*entity.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = FromUserEntity(user)
	}
	return responses
}
