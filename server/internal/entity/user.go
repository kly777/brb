package entity

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        uint      // 主键ID
	Username  string    // 用户名
	Password  string    // 密码（应加密存储）
	Role      Role      // 角色
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}
