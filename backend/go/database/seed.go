package database

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
)

// SeedData 数据库初始化种子数据
// 作者: cjx
// 邮箱: xx4125517@126.com
// 时间: 2025-12-23 14:50:00
func SeedData() error {
	logrus.Info("开始初始化种子数据...")

	// 1. 初始化默认集团公司
	company := model.Company{
		Name:   "集团总部",
		Code:   "HEADQUARTER",
		Type:   1,
		Status: 1,
	}
	// 使用 FirstOrCreate 简化逻辑，且因为是新库初始化，这里主要为了获取 ID
	if err := DB.Where("code = ?", company.Code).FirstOrCreate(&company).Error; err != nil {
		logrus.Errorf("初始化默认公司失败: %v", err)
		return err
	}
	logrus.Info("默认公司 '集团总部' 初始化成功")

	// 2. 初始化管理员角色
	adminRole := model.Role{
		Name:   "管理员",
		Code:   "admin",
		Status: 1,
	}
	if err := DB.Where("code = ?", adminRole.Code).FirstOrCreate(&adminRole).Error; err != nil {
		logrus.Errorf("初始化管理员角色失败: %v", err)
		return err
	}
	logrus.Info("管理员角色 'admin' 初始化成功")

	// 3. 初始化管理员账户
	adminUser := model.User{
		CompanyID: company.ID,
		Username:  "admin",
		Password:  "admin", // 配合当前系统的明文逻辑
		Nickname:  "系统管理员",
		Status:    1,
	}
	if err := DB.Where("username = ?", adminUser.Username).FirstOrCreate(&adminUser).Error; err != nil {
		logrus.Errorf("初始化管理员账号失败: %v", err)
		return err
	}

	// 绑定角色
	if err := DB.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
		logrus.Errorf("绑定管理员角色失败: %v", err)
		return err
	}
	logrus.Info("管理员账号 'admin' 初始化成功并绑定角色")

	// 4. 初始化数据字典
	statusEnabled := 1
	dictionaries := []model.DataDictionary{
		{Module: "system", Field: "status", Label: "启用", Value: "1", Sort: 1, Status: &statusEnabled},
		{Module: "system", Field: "status", Label: "禁用", Value: "0", Sort: 2, Status: &statusEnabled},
		{Module: "company", Field: "type", Label: "集团总部", Value: "1", Sort: 1, Status: &statusEnabled},
		{Module: "company", Field: "type", Label: "分子公司", Value: "2", Sort: 2, Status: &statusEnabled},
		{Module: "menu", Field: "type", Label: "菜单", Value: "1", Sort: 1, Status: &statusEnabled},
		{Module: "menu", Field: "type", Label: "按钮", Value: "2", Sort: 2, Status: &statusEnabled},
	}
	for _, dict := range dictionaries {
		if err := DB.Where("module = ? AND field = ? AND value = ?", dict.Module, dict.Field, dict.Value).FirstOrCreate(&dict).Error; err != nil {
			logrus.Errorf("初始化数据字典失败 [%s-%s]: %v", dict.Module, dict.Field, err)
			return err
		}
	}
	logrus.Info("数据字典初始化成功")

	// 5. 初始化系统菜单
	// 定义菜单结构
	type MenuNode struct {
		Name      string
		Path      string
		Component string
		Icon      string
		Sort      int
		Type      int
		Children  []MenuNode
	}

	systemMenus := []MenuNode{
		{
			Name: "系统管理", Path: "/system", Component: "Layout", Icon: "setting", Sort: 1, Type: 1,
			Children: []MenuNode{
				{Name: "用户管理", Path: "user", Component: "system/user/index", Icon: "user", Sort: 1, Type: 1},
				{Name: "角色管理", Path: "role", Component: "system/role/index", Icon: "peoples", Sort: 2, Type: 1},
				{Name: "菜单管理", Path: "menu", Component: "system/menu/index", Icon: "tree-table", Sort: 3, Type: 1},
				{Name: "公司管理", Path: "company", Component: "system/company/index", Icon: "tree", Sort: 4, Type: 1},
				{Name: "字段权限", Path: "field-permission", Component: "system/field-permission/index", Icon: "lock", Sort: 5, Type: 1},
				{Name: "数据字典", Path: "dict", Component: "system/dict/index", Icon: "list", Sort: 6, Type: 1},
			},
		},
		{
			Name: "业务功能", Path: "/business", Component: "Layout", Icon: "component", Sort: 2, Type: 1,
			Children: []MenuNode{
				{Name: "下载任务", Path: "download-task", Component: "business/download-task/index", Icon: "download", Sort: 1, Type: 1},
			},
		},
		{
			Name: "API测试", Path: "/api-test", Component: "Layout", Icon: "bug", Sort: 3, Type: 1,
			Children: []MenuNode{
				{Name: "测试用例", Path: "case", Component: "api-test/case/index", Icon: "file-text", Sort: 1, Type: 1},
				{Name: "测试历史", Path: "history", Component: "api-test/history/index", Icon: "history", Sort: 2, Type: 1},
			},
		},
	}

	// 递归创建菜单函数
	var createMenus func(menus []MenuNode, parentID uint) error
	createMenus = func(menus []MenuNode, parentID uint) error {
		for _, m := range menus {
			menu := model.Menu{
				ParentID:  parentID,
				Name:      m.Name,
				Path:      m.Path,
				Component: m.Component,
				Icon:      m.Icon,
				Sort:      m.Sort,
				Type:      m.Type,
				Status:    1,
			}
			// 根据路径和名称查重
			if err := DB.Where("name = ? AND parent_id = ?", m.Name, parentID).FirstOrCreate(&menu).Error; err != nil {
				return err
			}

			// 自动给管理员分配菜单权限
			if err := DB.Model(&adminRole).Association("Menus").Append(&menu); err != nil {
				logrus.Errorf("给管理员分配菜单权限失败 [%s]: %v", m.Name, err)
				// 继续执行，不阻断
			}

			// 递归创建子菜单
			if len(m.Children) > 0 {
				if err := createMenus(m.Children, menu.ID); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := createMenus(systemMenus, 0); err != nil {
		logrus.Errorf("初始化菜单失败: %v", err)
		return err
	}
	logrus.Info("系统菜单初始化成功，并已分配给管理员")

	logrus.Info("所有种子数据初始化完成")
	return nil
}
