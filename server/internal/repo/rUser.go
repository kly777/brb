package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type userRepo struct {
	base *BaseRepo[entity.User]
}

// NewUserRepo 创建新的用户Repository
func NewUserRepo(db *sql.DB) (*userRepo, error) {
	// 初始化数据库表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	baseRepo := NewBaseRepo[entity.User](db, "users")
	return &userRepo{base: baseRepo}, nil
}

// Create 创建新用户
func (r *userRepo) Create(user *entity.User) error {
	fields := map[string]any{
		"username":   user.Username,
		"password":   user.Password,
		"role":       string(user.Role),
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	// 如果ID已设置（用于更新），包含它，否则将自动生成
	if user.ID != 0 {
		fields["id"] = user.ID
	}

	_, err := r.base.Create(fields)
	return err
}

// GetByID 根据ID获取用户
func (r *userRepo) GetByID(id uint) (*entity.User, error) {
	query := "SELECT id, username, password, role, created_at, updated_at FROM users WHERE id = ?"
	row := r.base.db.QueryRow(query, id)

	user, err := r.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepo) GetByUsername(username string) (*entity.User, error) {
	query := "SELECT id, username, password, role, created_at, updated_at FROM users WHERE username = ?"
	row := r.base.db.QueryRow(query, username)

	user, err := r.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return user, nil
}

// GetAll 获取所有用户
func (r *userRepo) GetAll() ([]*entity.User, error) {
	query := "SELECT id, username, password, role, created_at, updated_at FROM users"
	rows, err := r.base.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

// Update 更新用户信息
func (r *userRepo) Update(user *entity.User) error {
	fields := map[string]any{
		"username":   user.Username,
		"password":   user.Password,
		"role":       string(user.Role),
		"updated_at": user.UpdatedAt,
	}

	return r.base.Update(user.ID, fields)
}

// Delete 删除用户
func (r *userRepo) Delete(id uint) error {
	return r.base.Delete(id)
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepo) ExistsByUsername(username string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	var count int
	err := r.base.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

// HaveID 检查用户ID是否存在
func (r *userRepo) HaveID(id uint) bool {
	query := "SELECT COUNT(*) FROM users WHERE id = ?"
	var count int
	err := r.base.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// scanUser 从数据库行扫描User实体
func (r *userRepo) scanUser(row any) (*entity.User, error) {
	var (
		id        uint
		username  string
		password  string
		role      string
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	var err error
	switch row := row.(type) {
	case *sql.Row:
		err = row.Scan(&id, &username, &password, &role, &createdAt, &updatedAt)
	case *sql.Rows:
		err = row.Scan(&id, &username, &password, &role, &createdAt, &updatedAt)
	default:
		return nil, fmt.Errorf("unsupported row type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	user := &entity.User{
		ID:       id,
		Username: username,
		Password: password,
		Role:     entity.Role(role),
	}

	// 处理时间字段
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}

	return user, nil
}
