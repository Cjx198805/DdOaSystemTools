package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	// 1. 获取关联的 API 配置
	var apiConfig model.APIConfig
	if err := s.db.First(&apiConfig, testCase.APIConfigID).Error; err != nil {
		logrus.Errorf("获取API配置失败: %v", err)
		return nil, err
	}

	// 2. 准备请求 URL
	baseURL := strings.TrimSuffix(apiConfig.BaseURL, "/")
	path := strings.TrimPrefix(apiConfig.Path, "/")
	fullURL := baseURL + "/" + path

	// 3. 准备请求头和参数（合并配置与用例）
	headers := make(map[string]string)
	if apiConfig.Headers != "" {
		json.Unmarshal([]byte(apiConfig.Headers), &headers)
	}
	if testCase.Headers != "" {
		tempHeaders := make(map[string]string)
		json.Unmarshal([]byte(testCase.Headers), &tempHeaders)
		for k, v := range tempHeaders {
			headers[k] = v
		}
	}

	params := make(map[string]interface{})
	if apiConfig.Params != "" {
		json.Unmarshal([]byte(apiConfig.Params), &params)
	}
	if testCase.Params != "" {
		tempParams := make(map[string]interface{})
		json.Unmarshal([]byte(testCase.Params), &tempParams)
		for k, v := range tempParams {
			params[k] = v
		}
	}

	// 4. 构建 HTTP 请求
	var bodyReader io.Reader
	method := strings.ToUpper(apiConfig.Method)
	if method == "" {
		method = "GET"
	}

	if method == "GET" || method == "DELETE" {
		if len(params) > 0 {
			if !strings.Contains(fullURL, "?") {
				fullURL += "?"
			} else {
				fullURL += "&"
			}
			for k, v := range params {
				fullURL += fmt.Sprintf("%s=%v&", k, v)
			}
			fullURL = strings.TrimSuffix(fullURL, "&")
		}
	} else {
		jsonBytes, _ := json.Marshal(params)
		bodyReader = strings.NewReader(string(jsonBytes))
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		logrus.Errorf("创建请求失败: %v", err)
		return nil, err
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if method != "GET" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 5. 执行请求并计时
	client := &http.Client{Timeout: 30 * time.Second}
	startTime := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(startTime).Milliseconds()

	// 6. 构造历史记录
	testHistory := &model.APITestHistory{
		CompanyID:    testCase.CompanyID,
		UserID:       userID,
		APIConfigID:  testCase.APIConfigID,
		TestCaseID:   testCase.ID,
		Name:         testCase.Name,
		Headers:      testCase.Headers,
		Params:       testCase.Params,
		ResponseTime: duration,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err != nil {
		testHistory.Status = "failed"
		testHistory.ErrorMessage = err.Error()
	} else {
		defer resp.Body.Close()
		testHistory.StatusCode = resp.StatusCode
		bodyBytes, _ := io.ReadAll(resp.Body)
		testHistory.ActualResult = string(bodyBytes)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			testHistory.Status = "success"
		} else {
			testHistory.Status = "failed"
			testHistory.ErrorMessage = fmt.Sprintf("HTTP Status %d", resp.StatusCode)
		}
	}

	// 7. 保存历史记录
	if err := s.db.Create(testHistory).Error; err != nil {
		logrus.Errorf("保存测试历史失败: %v", err)
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
