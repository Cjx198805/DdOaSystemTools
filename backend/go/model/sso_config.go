package model

import (
	"time"
)

// SSOConfig 身份验证（免登）配置模型
type SSOConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyID   uint      `gorm:"not null" json:"company_id"`
	AppID       string    `gorm:"size:100;not null" json:"app_id"`
	AppKey      string    `gorm:"size:100;not null" json:"app_key"`
	AppSecret   string    `gorm:"size:255;not null" json:"app_secret"`
	RedirectURL string    `gorm:"size:255" json:"redirect_url"`
	Scope       string    `gorm:"size:100" json:"scope"`
	Status      int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Company     Company    `gorm:"foreignKey:CompanyID" json:"company"`
}

// TableName 设置表名
func (SSOConfig) TableName() string {
	return "sso_config"
}
