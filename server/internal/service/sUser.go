package service

import (
	"brb/internal/entity"
	"brb/pkg/logger"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// userService 实现用户业务逻辑
type userService struct {
	userRepo userRepository
}

type userRepository interface {
	Create(user *entity.User) error
	GetByID(id uint) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
	ExistsByUsername(username string) (bool, error)
	HaveID(id uint) bool
}

// NewUserService 创建新的用户Service实例
func NewUserService(userRepo userRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *userService) Register(username, password string) (*entity.User, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &entity.User{
		Username:  username,
		Password:  string(hashedPassword),
		Role:      entity.RoleUser, // 注册时固定为user角色
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Tip.Printf("用户注册成功: %s (ID: %d)", username, user.ID)
	return user, nil
}

// Login 用户登录
func (s *userService) Login(username, password string) (*entity.User, error) {
	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	logger.Tip.Printf("用户登录成功: %s (ID: %d)", username, user.ID)
	return user, nil
}

// GetUserByID 根据ID获取用户信息
func (s *userService) GetUserByID(id uint) (*entity.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return user, nil
}

// GetAllUsers 获取所有用户（仅管理员可访问）
func (s *userService) GetAllUsers() ([]*entity.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}
	return users, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(id uint, username, password string, role entity.Role) (*entity.User, error) {
	// 获取现有用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 更新用户名（如果提供了新用户名）
	if username != "" {
		// 检查新用户名是否已被其他用户使用
		existingUser, err := s.userRepo.GetByUsername(username)
		if err == nil && existingUser.ID != id {
			return nil, fmt.Errorf("用户名已被其他用户使用")
		}
		user.Username = username
	}

	// 更新密码（如果提供了新密码）
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %w", err)
		}
		user.Password = string(hashedPassword)
	}

	// 更新角色（如果提供了新角色）
	if role != "" {
		user.Role = role
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	logger.Tip.Printf("用户信息已更新: %s (ID: %d)", user.Username, user.ID)
	return user, nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	// 检查用户是否存在
	if !s.userRepo.HaveID(id) {
		return fmt.Errorf("用户不存在")
	}

	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	logger.Tip.Printf("用户已删除: ID %d", id)
	return nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(id uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("密码更新失败: %w", err)
	}

	logger.Tip.Printf("用户密码已修改: %s (ID: %d)", user.Username, user.ID)
	return nil
}

// PromoteToAdmin 提升用户为管理员（仅管理员可操作）
func (s *userService) PromoteToAdmin(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	if user.Role == entity.RoleAdmin {
		return fmt.Errorf("用户已是管理员")
	}

	user.Role = entity.RoleAdmin
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("提升用户权限失败: %w", err)
	}

	logger.Tip.Printf("用户权限已提升为管理员: %s (ID: %d)", user.Username, user.ID)
	return nil
}

// DemoteToUser 降级用户为普通用户（仅管理员可操作）
func (s *userService) DemoteToUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	if user.Role == entity.RoleUser {
		return fmt.Errorf("用户已是普通用户")
	}

	user.Role = entity.RoleUser
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("降级用户权限失败: %w", err)
	}

	logger.Tip.Printf("用户权限已降级为普通用户: %s (ID: %d)", user.Username, user.ID)
	return nil
}