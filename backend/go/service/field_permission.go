package service

import (
	"errors"

	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// FieldPermissionService 字段权限服务
type FieldPermissionService struct{}

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
	if fieldPermission.Viewable == nil {
		fieldPermission.Viewable = util.IntPtr(1)
	}
	if fieldPermission.Editable == nil {
		fieldPermission.Editable = util.IntPtr(1)
	}
	if fieldPermission.ReportVisible == nil {
		fieldPermission.ReportVisible = util.IntPtr(1)
	}
	if fieldPermission.SpecialEdit == nil {
		fieldPermission.SpecialEdit = util.IntPtr(0)
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

// CheckFieldEditable 检查字段是否可编辑
// 核心逻辑：
// 1. 如果用户所属角色中有任意一个拥有 SpecialEdit = 1，则返回 true
// 2. 否则查看数据字典，如果字典设置为 Editable = 0，则返回 false
// 3. 默认返回 true
func (s *FieldPermissionService) CheckFieldEditable(userID uint, module, field string) (bool, error) {
	db := database.GetDB()

	// 1. 获取用户的角色ID列表
	var roleIDs []uint
	if err := db.Table("user_role").Where("user_id = ?", userID).Pluck("role_id", &roleIDs).Error; err != nil {
		logrus.Errorf("获取用户角色失败: %v", err)
		return false, err
	}

	if len(roleIDs) == 0 {
		return false, nil
	}

	// 2. 检查是否有特殊编辑权限 (SpecialEdit = 1)
	var specialCount int64
	err := db.Model(&model.FieldPermission{}).
		Where("role_id IN ? AND module = ? AND field = ? AND special_edit = 1", roleIDs, module, field).
		Count(&specialCount).Error
	if err != nil {
		logrus.Errorf("检查特殊编辑权限失败: %v", err)
		return false, err
	}

	if specialCount > 0 {
		return true, nil
	}

	// 3. 回退到数据字典检查
	var dict model.DataDictionary
	err = db.Where("module = ? AND field = ? AND status = 1", module, field).First(&dict).Error
	if err == nil {
		// 找到了字典项，检查是否可编辑
		return util.IntValue(dict.Editable, 1) == 1, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("查询数据字典失败: %v", err)
		return false, err
	}

	// 4. 既没有特殊权限也没有字典限制，默认允许编辑（或根据具体 RBAC 决定）
	return true, nil
}
