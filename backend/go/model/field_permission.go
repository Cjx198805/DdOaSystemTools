package model

import (
	"time"
)

// FieldPermission 字段权限模型
type FieldPermission struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	RoleID        uint       `gorm:"not null" json:"role_id"`
	Module        string     `gorm:"size:50;not null" json:"module"`  // 模块名称，如 user, company, api_config
	Field         string     `gorm:"size:50;not null" json:"field"`   // 字段名称
	Viewable      *int       `gorm:"default:1" json:"viewable"`       // 是否可查看 1:是 0:否
	Editable      *int       `gorm:"default:1" json:"editable"`       // 是否可编辑 1:是 0:否
	ReportVisible *int       `gorm:"default:1" json:"report_visible"` // 是否在报表中显示 1:是 0:否
	SpecialEdit   *int       `gorm:"default:0" json:"special_edit"`   // 特殊编辑权限 1:是 0:否，优先级高于数据字典
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Role Role `gorm:"foreignKey:RoleID" json:"role"`
}

// TableName 设置表名
func (FieldPermission) TableName() string {
	return "field_permission"
}
