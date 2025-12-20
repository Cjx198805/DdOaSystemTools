package service

import (
	"errors"
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// List 获取用户列表
func (s *UserService) List(page, pageSize int, companyID uint, username, nickname string, status int) ([]model.User, int64, error) {
	db := database.GetDB()
	
	var users []model.User
	var total int64
	
	// 构建查询
	query := db.Model(&model.User{})
	
	// 添加查询条件
	if companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+nickname+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		logrus.Errorf("获取用户总数失败: %v", err)
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Company").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		logrus.Errorf("获取用户列表失败: %v", err)
		return nil, 0, err
	}
	
	return users, total, nil
}

// Get 获取用户详情
func (s *UserService) Get(id uint) (*model.User, error) {
	db := database.GetDB()
	
	var user model.User
	if err := db.Preload("Company").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		logrus.Errorf("获取用户详情失败: %v", err)
		return nil, err
	}
	
	// 清空密码
	user.Password = ""
	
	return &user, nil
}

// Create 创建用户
func (s *UserService) Create(user *model.User) error {
	db := database.GetDB()
	
	// 检查公司是否存在
	var company model.Company
	if err := db.First(&company, user.CompanyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("公司不存在")
		}
		logrus.Errorf("检查公司是否存在失败: %v", err)
		return err
	}
	
	// 检查用户名是否已存在
	var existingUser model.User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查用户名是否存在失败: %v", err)
		return err
	}
	
	// 设置默认值
	if user.Status == 0 {
		user.Status = 1
	}
	
	// 创建用户
	if err := db.Create(user).Error; err != nil {
		logrus.Errorf("创建用户失败: %v", err)
		return err
	}
	
	return nil
}

// Update 更新用户
func (s *UserService) Update(user *model.User) error {
	db := database.GetDB()
	
	// 检查用户是否存在
	var existingUser model.User
	if err := db.First(&existingUser, user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		logrus.Errorf("获取用户失败: %v", err)
		return err
	}
	
	// 检查公司是否存在
	var company model.Company
	if err := db.First(&company, user.CompanyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("公司不存在")
		}
		logrus.Errorf("检查公司是否存在失败: %v", err)
		return err
	}
	
	// 检查用户名是否已被其他用户使用
	var duplicateUser model.User
	if err := db.Where("username = ? AND id != ?", user.Username, user.ID).First(&duplicateUser).Error; err == nil {
		return errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查用户名是否存在失败: %v", err)
		return err
	}
	
	// 不更新密码（密码更新通过单独的方法）
	user.Password = existingUser.Password
	
	// 更新用户
	if err := db.Save(user).Error; err != nil {
		logrus.Errorf("更新用户失败: %v", err)
		return err
	}
	
	return nil
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	db := database.GetDB()
	
	// 检查用户是否存在
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		logrus.Errorf("获取用户失败: %v", err)
		return err
	}
	
	// 软删除用户
	if err := db.Delete(&user).Error; err != nil {
		logrus.Errorf("删除用户失败: %v", err)
		return err
	}
	
	// 删除用户角色关联
	if err := db.Where("user_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
		logrus.Errorf("删除用户角色关联失败: %v", err)
		return err
	}
	
	return nil
}

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(id uint, newPassword string) error {
	db := database.GetDB()
	
	// 检查用户是否存在
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		logrus.Errorf("获取用户失败: %v", err)
		return err
	}
	
	// 更新密码
	user.Password = newPassword
	if err := db.Save(&user).Error; err != nil {
		logrus.Errorf("重置密码失败: %v", err)
		return err
	}
	
	return nil
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(id uint, oldPassword, newPassword string) error {
	db := database.GetDB()
	
	// 检查用户是否存在
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		logrus.Errorf("获取用户失败: %v", err)
		return err
	}
	
	// 检查旧密码是否正确
	if user.Password != oldPassword {
		return errors.New("旧密码不正确")
	}
	
	// 更新密码
	user.Password = newPassword
	if err := db.Save(&user).Error; err != nil {
		logrus.Errorf("更新密码失败: %v", err)
		return err
	}
	
	return nil
}

// GetRoles 获取用户角色
func (s *UserService) GetRoles(userID uint) ([]model.Role, error) {
	db := database.GetDB()
	
	var roles []model.Role
	if err := db.Table("role").Joins("JOIN user_role ON role.id = user_role.role_id").Where("user_role.user_id = ?", userID).Find(&roles).Error; err != nil {
		logrus.Errorf("获取用户角色失败: %v", err)
		return nil, err
	}
	
	return roles, nil
}

// AssignRoles 分配角色
func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	db := database.GetDB()
	
	// 检查用户是否存在
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		logrus.Errorf("获取用户失败: %v", err)
		return err
	}
	
	// 删除旧的角色关联
	if err := db.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		logrus.Errorf("删除旧角色关联失败: %v", err)
		return err
	}
	
	// 添加新的角色关联
	for _, roleID := range roleIDs {
		// 检查角色是否存在
		var role model.Role
		if err := db.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("角色不存在")
			}
			logrus.Errorf("获取角色失败: %v", err)
			return err
		}
		
		userRole := model.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := db.Create(&userRole).Error; err != nil {
			logrus.Errorf("创建用户角色关联失败: %v", err)
			return err
		}
	}
	
	return nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, error) {
	db := database.GetDB()
	
	// 查找用户
	var user model.User
	if err := db.Where("username = ? AND status = 1", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		logrus.Errorf("查找用户失败: %v", err)
		return nil, err
	}
	
	// 验证密码
	if user.Password != password {
		return nil, errors.New("用户名或密码错误")
	}
	
	// 清空密码
	user.Password = ""
	
	logrus.Infof("用户登录成功，用户名: %s", username)
	
	return &user, nil
}
