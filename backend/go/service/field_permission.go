package service

import (
	"errors"
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// FieldPermissionService 字段权限服务
type FieldPermissionService struct {}

// NewFieldPermissionService 创建字段权限服务实例
func NewFieldPermissionService() *FieldPermissionService {
	return &FieldPermissionService{}
}

// List 获取字段权限列表
func (s *FieldPermissionService) List(page, pageSize int, roleID uint, module, field string) ([]model.FieldPermission, int64, error) {
	db := database.GetDB()
	
	var fieldPermissions []model.FieldPermission
	var total int64
	
	// 构建查询
	query := db.Model(&model.FieldPermission{})
	
	// 添加查询条件
	if roleID > 0 {
		query = query.Where("role_id = ?", roleID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if field != "" {
		query = query.Where("field LIKE ?", "%"+field+"%")
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		logrus.Errorf("获取字段权限总数失败: %v", err)
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Role").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&fieldPermissions).Error; err != nil {
		logrus.Errorf("获取字段权限列表失败: %v", err)
		return nil, 0, err
	}
	
	return fieldPermissions, total, nil
}

// Get 获取字段权限详情
func (s *FieldPermissionService) Get(id uint) (*model.FieldPermission, error) {
	db := database.GetDB()
	
	var fieldPermission model.FieldPermission
	if err := db.Preload("Role").First(&fieldPermission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("字段权限不存在")
		}
		logrus.Errorf("获取字段权限详情失败: %v", err)
		return nil, err
	}
	
	return &fieldPermission, nil
}

// Create 创建字段权限
func (s *FieldPermissionService) Create(fieldPermission *model.FieldPermission) error {
	db := database.GetDB()
	
	// 检查角色是否存在
	var role model.Role
	if err := db.First(&role, fieldPermission.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		logrus.Errorf("检查角色是否存在失败: %v", err)
		return err
	}
	
	// 检查是否已存在相同的字段权限
	var existing model.FieldPermission
	result := db.Where("role_id = ? AND module = ? AND field = ?", fieldPermission.RoleID, fieldPermission.Module, fieldPermission.Field).First(&existing)
	if result.Error == nil {
		return errors.New("该角色的该字段权限已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查字段权限是否存在失败: %v", result.Error)
		return result.Error
	}
	
	// 设置默认值
	if fieldPermission.Viewable == 0 {
		fieldPermission.Viewable = 1
	}
	if fieldPermission.Editable == 0 {
		fieldPermission.Editable = 1
	}
	if fieldPermission.ReportVisible == 0 {
		fieldPermission.ReportVisible = 1
	}
	if fieldPermission.SpecialEdit == 0 {
		fieldPermission.SpecialEdit = 0
	}
	
	// 创建字段权限
	if err := db.Create(fieldPermission).Error; err != nil {
		logrus.Errorf("创建字段权限失败: %v", err)
		return err
	}
	
	return nil
}

// Update 更新字段权限
func (s *FieldPermissionService) Update(fieldPermission *model.FieldPermission) error {
	db := database.GetDB()
	
	// 检查字段权限是否存在
	var existing model.FieldPermission
	if err := db.First(&existing, fieldPermission.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("字段权限不存在")
		}
		logrus.Errorf("获取字段权限失败: %v", err)
		return err
	}
	
	// 检查角色是否存在
	var role model.Role
	if err := db.First(&role, fieldPermission.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		logrus.Errorf("检查角色是否存在失败: %v", err)
		return err
	}
	
	// 检查是否已存在相同的字段权限（排除当前记录）
	var duplicate model.FieldPermission
	result := db.Where("role_id = ? AND module = ? AND field = ? AND id != ?", fieldPermission.RoleID, fieldPermission.Module, fieldPermission.Field, fieldPermission.ID).First(&duplicate)
	if result.Error == nil {
		return errors.New("该角色的该字段权限已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查字段权限是否存在失败: %v", result.Error)
		return result.Error
	}
	
	// 更新字段权限
	if err := db.Save(fieldPermission).Error; err != nil {
		logrus.Errorf("更新字段权限失败: %v", err)
		return err
	}
	
	return nil
}

// Delete 删除字段权限
func (s *FieldPermissionService) Delete(id uint) error {
	db := database.GetDB()
	
	// 检查字段权限是否存在
	var fieldPermission model.FieldPermission
	if err := db.First(&fieldPermission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("字段权限不存在")
		}
		logrus.Errorf("获取字段权限失败: %v", err)
		return err
	}
	
	// 软删除字段权限
	if err := db.Delete(&fieldPermission).Error; err != nil {
		logrus.Errorf("删除字段权限失败: %v", err)
		return err
	}
	
	return nil
}

// GetByRoleAndModule 根据角色ID和模块获取字段权限
func (s *FieldPermissionService) GetByRoleAndModule(roleID uint, module string) ([]model.FieldPermission, error) {
	db := database.GetDB()
	
	var fieldPermissions []model.FieldPermission
	if err := db.Where("role_id = ? AND module = ?", roleID, module).Find(&fieldPermissions).Error; err != nil {
		logrus.Errorf("根据角色ID和模块获取字段权限失败: %v", err)
		return nil, err
	}
	
	return fieldPermissions, nil
}

// GetByRoleIDsAndModule 根据角色ID列表和模块获取字段权限
func (s *FieldPermissionService) GetByRoleIDsAndModule(roleIDs []uint, module string) ([]model.FieldPermission, error) {
	db := database.GetDB()
	
	var fieldPermissions []model.FieldPermission
	if err := db.Where("role_id IN ? AND module = ?", roleIDs, module).Find(&fieldPermissions).Error; err != nil {
		logrus.Errorf("根据角色ID列表和模块获取字段权限失败: %v", err)
		return nil, err
	}
	
	return fieldPermissions, nil
}
