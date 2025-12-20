package model

import (
	"time"
)

// Menu 菜单模型
type Menu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"default:0" json:"parent_id"` // 父菜单ID，0表示一级菜单
	Name      string    `gorm:"size:50;not null" json:"name"`
	Path      string    `gorm:"size:100" json:"path"`
	Component string    `gorm:"size:100" json:"component"`
	Icon      string    `gorm:"size:50" json:"icon"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Type      int       `gorm:"default:1" json:"type"` // 1: 菜单, 2: 按钮
	Status    int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 设置表名
func (Menu) TableName() string {
	return "menu"
}

// RoleMenu 角色菜单关联模型
type RoleMenu struct {
	RoleID uint `gorm:"primaryKey" json:"role_id"`
	MenuID uint `gorm:"primaryKey" json:"menu_id"`
}

// TableName 设置表名
func (RoleMenu) TableName() string {
	return "role_menu"
}
