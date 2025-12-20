package service

import (
	"errors"
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct {}

// NewRoleService 创建角色服务实例
func NewRoleService() *RoleService {
	return &RoleService{}
}

// List 获取角色列表
func (s *RoleService) List(page, pageSize int, name, code string, status int) ([]model.Role, int64, error) {
	db := database.GetDB()
	
	var roles []model.Role
	var total int64
	
	// 构建查询
	query := db.Model(&model.Role{})
	
	// 添加查询条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		logrus.Errorf("获取角色总数失败: %v", err)
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&roles).Error; err != nil {
		logrus.Errorf("获取角色列表失败: %v", err)
		return nil, 0, err
	}
	
	return roles, total, nil
}

// Get 获取角色详情
func (s *RoleService) Get(id uint) (*model.Role, error) {
	db := database.GetDB()
	
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		logrus.Errorf("获取角色详情失败: %v", err)
		return nil, err
	}
	
	return &role, nil
}

// Create 创建角色
func (s *RoleService) Create(role *model.Role) error {
	db := database.GetDB()
	
	// 检查角色编码是否已存在
	var existingRole model.Role
	if err := db.Where("code = ?", role.Code).First(&existingRole).Error; err == nil {
		return errors.New("角色编码已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查角色编码是否存在失败: %v", err)
		return err
	}
	
	// 设置默认值
	if role.Status == 0 {
		role.Status = 1
	}
	
	// 创建角色
	if err := db.Create(role).Error; err != nil {
		logrus.Errorf("创建角色失败: %v", err)
		return err
	}
	
	return nil
}

// Update 更新角色
func (s *RoleService) Update(role *model.Role) error {
	db := database.GetDB()
	
	// 检查角色是否存在
	var existingRole model.Role
	if err := db.First(&existingRole, role.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		logrus.Errorf("获取角色失败: %v", err)
		return err
	}
	
	// 检查角色编码是否已被其他角色使用
	var duplicateRole model.Role
	if err := db.Where("code = ? AND id != ?", role.Code, role.ID).First(&duplicateRole).Error; err == nil {
		return errors.New("角色编码已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查角色编码是否存在失败: %v", err)
		return err
	}
	
	// 更新角色
	if err := db.Save(role).Error; err != nil {
		logrus.Errorf("更新角色失败: %v", err)
		return err
	}
	
	return nil
}

// Delete 删除角色
func (s *RoleService) Delete(id uint) error {
	db := database.GetDB()
	
	// 检查角色是否存在
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		logrus.Errorf("获取角色失败: %v", err)
		return err
	}
	
	// 软删除角色
	if err := db.Delete(&role).Error; err != nil {
		logrus.Errorf("删除角色失败: %v", err)
		return err
	}
	
	// 删除角色菜单关联
	if err := db.Where("role_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
		logrus.Errorf("删除角色菜单关联失败: %v", err)
		return err
	}
	
	// 删除用户角色关联
	if err := db.Where("role_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
		logrus.Errorf("删除用户角色关联失败: %v", err)
		return err
	}
	
	return nil
}

// GetMenus 获取角色菜单
func (s *RoleService) GetMenus(roleID uint) ([]model.Menu, error) {
	db := database.GetDB()
	
	var menus []model.Menu
	if err := db.Table("menu").Joins("JOIN role_menu ON menu.id = role_menu.menu_id").Where("role_menu.role_id = ? AND menu.status = 1", roleID).Order("menu.sort ASC").Find(&menus).Error; err != nil {
		logrus.Errorf("获取角色菜单失败: %v", err)
		return nil, err
	}
	
	return menus, nil
}

// AssignMenus 分配菜单
func (s *RoleService) AssignMenus(roleID uint, menuIDs []uint) error {
	db := database.GetDB()
	
	// 检查角色是否存在
	var role model.Role
	if err := db.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		logrus.Errorf("获取角色失败: %v", err)
		return err
	}
	
	// 删除旧的菜单关联
	if err := db.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		logrus.Errorf("删除旧菜单关联失败: %v", err)
		return err
	}
	
	// 添加新的菜单关联
	for _, menuID := range menuIDs {
		// 检查菜单是否存在
		var menu model.Menu
		if err := db.First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("菜单不存在")
			}
			logrus.Errorf("获取菜单失败: %v", err)
			return err
		}
		
		roleMenu := model.RoleMenu{
			RoleID: roleID,
			MenuID: menuID,
		}
		if err := db.Create(&roleMenu).Error; err != nil {
			logrus.Errorf("创建角色菜单关联失败: %v", err)
			return err
		}
	}
	
	return nil
}
