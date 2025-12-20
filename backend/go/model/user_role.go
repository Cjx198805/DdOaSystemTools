package model

import (
	"time"
)

// UserRole 用户角色关联表
// 用于实现单用户多角色分配
type UserRole struct {
	CompanyID uint      `gorm:"not null;index:idx_company_user_role" json:"company_id"`
	UserID    uint      `gorm:"primaryKey;not null;index:idx_user_role" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey;not null;index:idx_user_role" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName 设置表名
func (UserRole) TableName() string {
	return "user_role"
}