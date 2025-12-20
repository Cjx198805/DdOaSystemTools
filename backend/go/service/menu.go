package service

import (
	"errors"
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// MenuService 菜单服务
type MenuService struct {}

// NewMenuService 创建菜单服务实例
func NewMenuService() *MenuService {
	return &MenuService{}
}

// List 获取菜单列表
func (s *MenuService) List(page, pageSize int, name, path string, status int) ([]model.Menu, int64, error) {
	db := database.GetDB()
	
	var menus []model.Menu
	var total int64
	
	// 构建查询
	query := db.Model(&model.Menu{})
	
	// 添加查询条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if path != "" {
		query = query.Where("path LIKE ?", "%"+path+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		logrus.Errorf("获取菜单总数失败: %v", err)
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		logrus.Errorf("获取菜单列表失败: %v", err)
		return nil, 0, err
	}
	
	return menus, total, nil
}

// Get 获取菜单详情
func (s *MenuService) Get(id uint) (*model.Menu, error) {
	db := database.GetDB()
	
	var menu model.Menu
	if err := db.First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("菜单不存在")
		}
		logrus.Errorf("获取菜单详情失败: %v", err)
		return nil, err
	}
	
	return &menu, nil
}

// Create 创建菜单
func (s *MenuService) Create(menu *model.Menu) error {
	db := database.GetDB()
	
	// 设置默认值
	if menu.ParentID == 0 {
		menu.ParentID = 0
	}
	if menu.Sort == 0 {
		menu.Sort = 0
	}
	if menu.Type == 0 {
		menu.Type = 1
	}
	if menu.Status == 0 {
		menu.Status = 1
	}
	
	// 创建菜单
	if err := db.Create(menu).Error; err != nil {
		logrus.Errorf("创建菜单失败: %v", err)
		return err
	}
	
	return nil
}

// Update 更新菜单
func (s *MenuService) Update(menu *model.Menu) error {
	db := database.GetDB()
	
	// 检查菜单是否存在
	var existingMenu model.Menu
	if err := db.First(&existingMenu, menu.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		logrus.Errorf("获取菜单失败: %v", err)
		return err
	}
	
	// 不能将自己设为自己的父菜单
	if menu.ID == menu.ParentID {
		return errors.New("不能将自己设为自己的父菜单")
	}
	
	// 更新菜单
	if err := db.Save(menu).Error; err != nil {
		logrus.Errorf("更新菜单失败: %v", err)
		return err
	}
	
	return nil
}

// Delete 删除菜单
func (s *MenuService) Delete(id uint) error {
	db := database.GetDB()
	
	// 检查菜单是否存在
	var menu model.Menu
	if err := db.First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		logrus.Errorf("获取菜单失败: %v", err)
		return err
	}
	
	// 检查是否有子菜单
	var childCount int64
	if err := db.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		logrus.Errorf("检查子菜单失败: %v", err)
		return err
	}
	
	if childCount > 0 {
		return errors.New("该菜单下存在子菜单，无法删除")
	}
	
	// 删除菜单
	if err := db.Delete(&menu).Error; err != nil {
		logrus.Errorf("删除菜单失败: %v", err)
		return err
	}
	
	// 删除角色菜单关联
	if err := db.Where("menu_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
		logrus.Errorf("删除角色菜单关联失败: %v", err)
		return err
	}
	
	return nil
}

// GetTree 获取菜单树形结构
func (s *MenuService) GetTree() ([]model.Menu, error) {
	db := database.GetDB()
	
	// 获取所有启用的菜单
	var menus []model.Menu
	if err := db.Where("status = 1").Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		logrus.Errorf("获取菜单列表失败: %v", err)
		return nil, err
	}
	
	// 构建树形结构
	tree := buildMenuTree(menus, 0)
	
	return tree, nil
}

// buildMenuTree 构建菜单树形结构
func buildMenuTree(menus []model.Menu, parentID uint) []model.Menu {
	var tree []model.Menu
	
	for _, menu := range menus {
		if menu.ParentID == parentID {
			// 直接添加当前菜单，在前端构建树形结构
			tree = append(tree, menu)
		}
	}
	
	return tree
}

// GetByParentID 根据父ID获取子菜单
func (s *MenuService) GetByParentID(parentID uint) ([]model.Menu, error) {
	db := database.GetDB()
	
	var menus []model.Menu
	if err := db.Where("parent_id = ? AND status = 1", parentID).Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		logrus.Errorf("获取子菜单失败: %v", err)
		return nil, err
	}
	
	return menus, nil
}

// GetAll 获取所有菜单
func (s *MenuService) GetAll(status int) ([]model.Menu, error) {
	db := database.GetDB()
	
	var menus []model.Menu
	
	// 构建查询
	query := db.Model(&model.Menu{})
	
	// 添加查询条件
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		logrus.Errorf("获取所有菜单失败: %v", err)
		return nil, err
	}
	
	return menus, nil
}
