package model

import (
	"time"
)

// Log 日志模型
type Log struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID uint      `gorm:"not null" json:"company_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	Type      int       `gorm:"default:1" json:"type"` // 1: 操作日志, 2: 系统日志, 3: API 调用日志
	Action    string    `gorm:"size:100" json:"action"`
	Content   string    `gorm:"type:text" json:"content"`
	IP        string    `gorm:"size:50" json:"ip"`
	Status    int       `gorm:"default:1" json:"status"` // 1: 成功, 0: 失败
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Company   Company    `gorm:"foreignKey:CompanyID" json:"company"`
}

// TableName 设置表名
func (Log) TableName() string {
	return "log"
}
