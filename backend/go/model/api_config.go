package model

import (
	"time"
)

// APIConfig API 配置模型
type APIConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyID   uint      `gorm:"not null" json:"company_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Code        string    `gorm:"size:100;uniqueIndex" json:"code"`
	Version     string    `gorm:"size:20;not null" json:"version"` // API 版本，如 v1, v2
	Type        int       `gorm:"default:1" json:"type"` // 1: 新版 API, 2: 旧版 API
	BaseURL     string    `gorm:"size:200" json:"base_url"`
	Path        string    `gorm:"size:200" json:"path"`
	Method      string    `gorm:"size:10" json:"method"` // GET, POST, PUT, DELETE
	Params      string    `gorm:"type:text" json:"params"` // 参数配置，JSON 格式
	Headers     string    `gorm:"type:text" json:"headers"` // 请求头配置，JSON 格式
	Description string    `gorm:"type:text" json:"description"`
	Status      int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Company     Company    `gorm:"foreignKey:CompanyID" json:"company"`
}

// TableName 设置表名
func (APIConfig) TableName() string {
	return "api_config"
}
