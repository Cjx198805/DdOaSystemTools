package model

import (
	"time"
)

// AccessToken accessToken 模型
type AccessToken struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyID   uint      `gorm:"not null;uniqueIndex" json:"company_id"`
	AppKey      string    `gorm:"size:100;not null" json:"app_key"`
	AppSecret   string    `gorm:"size:255;not null" json:"app_secret"`
	AccessToken string    `gorm:"size:500" json:"access_token"`
	ExpiresIn   int       `gorm:"default:7200" json:"expires_in"` // 过期时间（秒）
	ExpiresAt   time.Time `json:"expires_at"`                     // 过期时间
	RefreshAt   time.Time `json:"refresh_at"`                     // 刷新时间
	Status      int       `gorm:"default:1" json:"status"`       // 1: 有效, 0: 无效
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Company     Company    `gorm:"foreignKey:CompanyID" json:"company"`
}

// TableName 设置表名
func (AccessToken) TableName() string {
	return "access_token"
}
