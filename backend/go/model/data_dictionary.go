package model

import (
	"time"
)

// DataDictionary 数据字典模型
type DataDictionary struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Module       string     `gorm:"size:50;not null;index:idx_module_field" json:"module"` // 模块
	Field        string     `gorm:"size:50;not null;index:idx_module_field" json:"field"`  // 字段
	Label        string     `gorm:"size:100;not null" json:"label"`                        // 显示标签
	Value        string     `gorm:"size:100" json:"value"`                                 // 存储值 (新增)
	Sort         int        `gorm:"default:0" json:"sort"`                                 // 排序 (新增)
	Description  string     `gorm:"type:text" json:"description"`                          // 字段描述 (恢复)
	Required     *int       `gorm:"default:0" json:"required"`                             // 是否必填 (恢复)
	Editable     *int       `gorm:"default:1" json:"editable"`                             // 是否可编辑 (恢复)
	DataType     string     `gorm:"size:20;default:'string'" json:"data_type"`             // 数据类型 (恢复)
	MaxLength    int        `gorm:"default:0" json:"max_length"`                           // 最大长度 (恢复)
	Options      string     `gorm:"type:text" json:"options"`                              // 可选值 (恢复)
	DefaultValue string     `gorm:"size:255" json:"default_value"`                         // 默认值 (恢复)
	Status       *int       `gorm:"default:1" json:"status"`                               // 状态
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 设置表名
func (DataDictionary) TableName() string {
	return "data_dictionary"
}
