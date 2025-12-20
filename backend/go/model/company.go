package model

import (
	"time"
)

// Company 集团公司模型
type Company struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"default:0" json:"parent_id"` // 父公司ID，0表示集团总部
	Name      string    `gorm:"size:100;not null" json:"name"`
	Code      string    `gorm:"size:50;uniqueIndex" json:"code"`
	Type      int       `gorm:"default:1" json:"type"` // 1: 集团总部, 2: 分子公司
	Status    int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 设置表名
func (Company) TableName() string {
	return "company"
}
