package service

import (
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"gorm.io/gorm"
	"time"
)

// APITestService API测试服务
type APITestService struct {
	db *gorm.DB
}

// NewAPITestService 创建API测试服务实例
func NewAPITestService() *APITestService {
	return &APITestService{
		db: database.GetDB(),
	}
}

// ListTestCases 获取API测试用例列表
func (s *APITestService) ListTestCases(page, pageSize int, companyID, apiConfigID uint, name string, status int) ([]model.APITestCase, int64, error) {
	var testCases []model.APITestCase
	var total int64

	query := s.db.Model(&model.APITestCase{})

	// 添加查询条件
	if companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}

	if apiConfigID > 0 {
		query = query.Where("api_config_id = ?", apiConfigID)
	}

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if status != -1 {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&testCases).Error; err != nil {
		return nil, 0, err
	}

	return testCases, total, nil
}

// GetTestCase 获取API测试用例详情
func (s *APITestService) GetTestCase(id uint) (*model.APITestCase, error) {
	var testCase model.APITestCase
	if err := s.db.First(&testCase, id).Error; err != nil {
		return nil, err
	}
	return &testCase, nil
}

// CreateTestCase 创建API测试用例
func (s *APITestService) CreateTestCase(testCase *model.APITestCase) error {
	return s.db.Create(testCase).Error
}

// UpdateTestCase 更新API测试用例
func (s *APITestService) UpdateTestCase(testCase *model.APITestCase) error {
	return s.db.Save(testCase).Error
}

// DeleteTestCase 删除API测试用例
func (s *APITestService) DeleteTestCase(id uint) error {
	return s.db.Delete(&model.APITestCase{}, id).Error
}

// RunTestCase 执行API测试用例
func (s *APITestService) RunTestCase(userID uint, testCase *model.APITestCase) (*model.APITestHistory, error) {
	// 创建测试历史记录
	testHistory := &model.APITestHistory{
		CompanyID:   testCase.CompanyID,
		UserID:      userID,
		APIConfigID: testCase.APIConfigID,
		TestCaseID:  testCase.ID,
		Name:        testCase.Name,
		Headers:     testCase.Headers,
		Params:      testCase.Params,
		Status:      "running",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 这里应该执行实际的API测试，暂时模拟成功结果
	testHistory.Status = "success"
	testHistory.ActualResult = "{\"code\": 200, \"message\": \"测试成功\", \"data\": null}"
	testHistory.ResponseTime = 123
	testHistory.StatusCode = 200

	// 保存测试历史记录
	if err := s.db.Create(testHistory).Error; err != nil {
		return nil, err
	}

	return testHistory, nil
}

// ListTestHistory 获取API测试历史记录列表
func (s *APITestService) ListTestHistory(page, pageSize int, companyID, userID, apiConfigID, testCaseID uint, name, status string) ([]model.APITestHistory, int64, error) {
	var testHistories []model.APITestHistory
	var total int64

	query := s.db.Model(&model.APITestHistory{})

	// 添加查询条件
	if companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if apiConfigID > 0 {
		query = query.Where("api_config_id = ?", apiConfigID)
	}

	if testCaseID > 0 {
		query = query.Where("api_test_id = ?", testCaseID)
	}

	if name != "" {
		query = query.Where("api_name LIKE ?", "%"+name+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&testHistories).Error; err != nil {
		return nil, 0, err
	}

	return testHistories, total, nil
}

// GetTestHistory 获取API测试历史记录详情
func (s *APITestService) GetTestHistory(id uint) (*model.APITestHistory, error) {
	var testHistory model.APITestHistory
	if err := s.db.First(&testHistory, id).Error; err != nil {
		return nil, err
	}
	return &testHistory, nil
}

// DeleteTestHistory 删除API测试历史记录
func (s *APITestService) DeleteTestHistory(id uint) error {
	return s.db.Delete(&model.APITestHistory{}, id).Error
}

// ClearTestHistory 清空API测试历史记录
func (s *APITestService) ClearTestHistory(companyID uint) error {
	return s.db.Where("company_id = ?", companyID).Delete(&model.APITestHistory{}).Error
}