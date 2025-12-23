package model

import (
	"time"
)

// DataDictionary 数据字典模型
type DataDictionary struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Module       string     `gorm:"size:50;not null" json:"module"`                     // 模块名称，如 user, company, api_config
	Field        string     `gorm:"size:50;not null" json:"field"`                      // 字段名称
	Label        string     `gorm:"size:100;not null" json:"label"`                     // 字段显示名称
	Description  string     `gorm:"type:text" json:"description"`                       // 字段描述
	Required     *int       `gorm:"default:0" json:"required"`                          // 是否必填 1:是 0:否
	Editable     *int       `gorm:"default:1" json:"editable"`                          // 创建后是否可编辑 1:是 0:否
	DataType     string     `gorm:"size:20;not null;default:'string'" json:"data_type"` // 数据类型，如 string, number, bool, date
	MaxLength    int        `gorm:"default:0" json:"max_length"`                        // 最大长度 0表示无限制
	Options      string     `gorm:"type:text" json:"options"`                           // 可选值，JSON格式，如 {"options":["option1","option2"]}
	DefaultValue string     `gorm:"size:255" json:"default_value"`                      // 默认值
	Status       *int       `gorm:"default:1" json:"status"`                            // 状态 1:启用 0:禁用
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 设置表名
func (DataDictionary) TableName() string {
	return "data_dictionary"
}
