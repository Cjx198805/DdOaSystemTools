package model

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:50;not null" json:"name"`
	Code      string     `gorm:"size:50;uniqueIndex:uni_role_code" json:"code"`
	Status    int        `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Users     []User     `gorm:"many2many:user_roles;" json:"users,omitempty"`
	UserRoles []UserRole `gorm:"foreignKey:RoleID" json:"user_roles,omitempty"`
	Menus     []Menu     `gorm:"many2many:role_menu;" json:"menus,omitempty"`
}

// TableName 设置表名
func (Role) TableName() string {
	return "roles"
}
