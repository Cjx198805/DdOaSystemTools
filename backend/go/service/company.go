package service

import (
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
)

// CompanyService 集团公司服务
type CompanyService struct {}

// NewCompanyService 创建集团公司服务
func NewCompanyService() *CompanyService {
	return &CompanyService{}
}

// List 获取集团公司列表
func (s *CompanyService) List(page, pageSize int, name, code string, companyType int) ([]model.Company, int64, error) {
	db := database.GetDB()
	
	var companies []model.Company
	var total int64
	
	// 构建查询
	query := db.Model(&model.Company{})
	
	// 添加查询条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if companyType > 0 {
		query = query.Where("type = ?", companyType)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		logrus.Errorf("获取集团公司总数失败: %v", err)
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("parent_id ASC, sort ASC, id ASC").Find(&companies).Error; err != nil {
		logrus.Errorf("获取集团公司列表失败: %v", err)
		return nil, 0, err
	}
	
	return companies, total, nil
}

// Create 创建集团公司
func (s *CompanyService) Create(company *model.Company) error {
	db := database.GetDB()
	
	if err := db.Create(company).Error; err != nil {
		logrus.Errorf("创建集团公司失败: %v", err)
		return err
	}
	
	return nil
}

// Get 获取集团公司详情
func (s *CompanyService) Get(id uint) (*model.Company, error) {
	db := database.GetDB()
	
	var company model.Company
	if err := db.First(&company, id).Error; err != nil {
		logrus.Errorf("获取集团公司详情失败: %v", err)
		return nil, err
	}
	
	return &company, nil
}

// Update 更新集团公司
func (s *CompanyService) Update(company *model.Company) error {
	db := database.GetDB()
	
	if err := db.Save(company).Error; err != nil {
		logrus.Errorf("更新集团公司失败: %v", err)
		return err
	}
	
	return nil
}

// Delete 删除集团公司
func (s *CompanyService) Delete(id uint) error {
	db := database.GetDB()
	
	// 检查是否有子公司
	var childCount int64
	if err := db.Model(&model.Company{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		logrus.Errorf("检查子公司失败: %v", err)
		return err
	}
	
	if childCount > 0 {
		logrus.Errorf("删除集团公司失败，存在 %d 个子公司", childCount)
		return nil
	}
	
	if err := db.Delete(&model.Company{}, id).Error; err != nil {
		logrus.Errorf("删除集团公司失败: %v", err)
		return err
	}
	
	return nil
}

// GetTree 获取集团公司树形结构
func (s *CompanyService) GetTree() ([]model.Company, error) {
	db := database.GetDB()
	
	var companies []model.Company
	if err := db.Order("parent_id ASC, sort ASC, id ASC").Find(&companies).Error; err != nil {
		logrus.Errorf("获取集团公司列表失败: %v", err)
		return nil, err
	}
	
	// 构建树形结构
	// TODO: 实现树形结构构建逻辑
	// 这里暂时返回列表，后续优化为树形结构
	
	return companies, nil
}
