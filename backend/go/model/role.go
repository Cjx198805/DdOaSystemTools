package model

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Code      string    `gorm:"size:50;uniqueIndex" json:"code"`
	Status    int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 设置表名
func (Role) TableName() string {
	return "role"
}

// UserRole 用户角色关联模型
type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	RoleID uint `gorm:"primaryKey" json:"role_id"`
}

// TableName 设置表名
func (UserRole) TableName() string {
	return "user_role"
}
