package model

// UserRole 用户角色关联表
// 用于实现单用户多角色分配
type UserRole struct {
	UserID uint `gorm:"primaryKey;not null" json:"user_id"`
	RoleID uint `gorm:"primaryKey;not null" json:"role_id"`
}

// TableName 设置表名
func (UserRole) TableName() string {
	return "user_roles"
}
