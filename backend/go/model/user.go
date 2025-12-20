package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID uint      `gorm:"not null" json:"company_id"`
	Username  string    `gorm:"size:50;uniqueIndex" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"password,omitempty"`
	Nickname  string    `gorm:"size:100" json:"nickname"`
	Email     string    `gorm:"size:100" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Status    int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Company   Company    `gorm:"foreignKey:CompanyID" json:"company"`
}

// TableName 设置表名
func (User) TableName() string {
	return "user"
}
