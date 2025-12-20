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
	
	// 关联关系
	Users     []User     `gorm:"many2many:user_role;" json:"users,omitempty"`
	UserRoles []UserRole `gorm:"foreignKey:RoleID" json:"user_roles,omitempty"`
}

// TableName 设置表名
func (Role) TableName() string {
	return "role"
}
