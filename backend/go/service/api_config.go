package service

import (
	"encoding/json"
	"errors"
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

// APIConfigService API配置服务
type APIConfigService struct {}

// NewAPIConfigService 创建API配置服务实例
func NewAPIConfigService() *APIConfigService {
	return &APIConfigService{}
}

// List 获取API配置列表
func (s *APIConfigService) List(page, pageSize int, companyID uint, name, code string, status int) ([]model.APIConfig, int64, error) {
	db := database.GetDB()
	
	var apiConfigs []model.APIConfig
	var total int64
	
	// 构建查询
	query := db.Model(&model.APIConfig{})
	
	// 添加查询条件
	if companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}
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
		logrus.Errorf("获取API配置总数失败: %v", err)
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Company").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&apiConfigs).Error; err != nil {
		logrus.Errorf("获取API配置列表失败: %v", err)
		return nil, 0, err
	}
	
	return apiConfigs, total, nil
}

// Get 获取API配置详情
func (s *APIConfigService) Get(id uint) (*model.APIConfig, error) {
	db := database.GetDB()
	
	var apiConfig model.APIConfig
	if err := db.Preload("Company").First(&apiConfig, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("API配置不存在")
		}
		logrus.Errorf("获取API配置详情失败: %v", err)
		return nil, err
	}
	
	return &apiConfig, nil
}

// Create 创建API配置
func (s *APIConfigService) Create(apiConfig *model.APIConfig) error {
	db := database.GetDB()
	
	// 检查公司是否存在
	var company model.Company
	if err := db.First(&company, apiConfig.CompanyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("公司不存在")
		}
		logrus.Errorf("检查公司是否存在失败: %v", err)
		return err
	}
	
	// 检查code是否已存在
	var existing model.APIConfig
	if err := db.Where("code = ?", apiConfig.Code).First(&existing).Error; err == nil {
		return errors.New("API配置编码已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查API配置编码是否存在失败: %v", err)
		return err
	}
	
	// 设置默认值
	if apiConfig.Status == 0 {
		apiConfig.Status = 1
	}
	
	// 创建API配置
	if err := db.Create(apiConfig).Error; err != nil {
		logrus.Errorf("创建API配置失败: %v", err)
		return err
	}
	
	return nil
}

// Update 更新API配置
func (s *APIConfigService) Update(apiConfig *model.APIConfig) error {
	db := database.GetDB()
	
	// 检查是否存在
	var existing model.APIConfig
	if err := db.First(&existing, apiConfig.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("API配置不存在")
		}
		logrus.Errorf("获取API配置失败: %v", err)
		return err
	}
	
	// 检查公司是否存在
	var company model.Company
	if err := db.First(&company, apiConfig.CompanyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("公司不存在")
		}
		logrus.Errorf("检查公司是否存在失败: %v", err)
		return err
	}
	
	// 检查code是否已被其他记录使用
	var duplicate model.APIConfig
	if err := db.Where("code = ? AND id != ?", apiConfig.Code, apiConfig.ID).First(&duplicate).Error; err == nil {
		return errors.New("API配置编码已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查API配置编码是否存在失败: %v", err)
		return err
	}
	
	// 更新API配置
	if err := db.Save(apiConfig).Error; err != nil {
		logrus.Errorf("更新API配置失败: %v", err)
		return err
	}
	
	return nil
}

// Delete 删除API配置
func (s *APIConfigService) Delete(id uint) error {
	db := database.GetDB()
	
	// 检查是否存在
	var apiConfig model.APIConfig
	if err := db.First(&apiConfig, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("API配置不存在")
		}
		logrus.Errorf("获取API配置失败: %v", err)
		return err
	}
	
	// 删除API配置
	if err := db.Delete(&apiConfig).Error; err != nil {
		logrus.Errorf("删除API配置失败: %v", err)
		return err
	}
	
	return nil
}

// Test 测试API配置
func (s *APIConfigService) Test(apiConfig *model.APIConfig) (map[string]interface{}, error) {
	// 解析参数
	var params map[string]interface{}
	if apiConfig.Params != "" {
		if err := json.Unmarshal([]byte(apiConfig.Params), &params); err != nil {
			return nil, errors.New(fmt.Sprintf("参数格式错误: %v", err))
		}
	}
	
	// 解析请求头
	var headers map[string]string
	if apiConfig.Headers != "" {
		if err := json.Unmarshal([]byte(apiConfig.Headers), &headers); err != nil {
			return nil, errors.New(fmt.Sprintf("请求头格式错误: %v", err))
		}
	}
	
	// 构建请求URL
	url := apiConfig.BaseURL
	if !strings.HasSuffix(url, "/") && !strings.HasPrefix(apiConfig.Path, "/") {
		url += "/"
	}
	url += apiConfig.Path
	
	// 如果是GET请求，将参数添加到URL
	if apiConfig.Method == "GET" && len(params) > 0 {
		url += "?"
		for k, v := range params {
			url += fmt.Sprintf("%s=%v&", k, v)
		}
		url = strings.TrimSuffix(url, "&")
	}
	
	// 创建请求
	var req *http.Request
	var err error
	
	// 记录开始时间
	startTime := time.Now()
	
	if apiConfig.Method == "GET" {
		req, err = http.NewRequest(apiConfig.Method, url, nil)
	} else {
		// 对于非GET请求，将参数转换为JSON
		jsonParams, err := json.Marshal(params)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("参数转换为JSON失败: %v", err))
		}
		req, err = http.NewRequest(apiConfig.Method, url, strings.NewReader(string(jsonParams)))
	}
	
	if err != nil {
		return nil, errors.New(fmt.Sprintf("创建请求失败: %v", err))
	}
	
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	
	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("发送请求失败: %v", err))
	}
	defer resp.Body.Close()
	
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("读取响应失败: %v", err))
	}
	
	// 解析响应
	var respData map[string]interface{}
	if err := json.Unmarshal(respBody, &respData); err != nil {
		// 如果响应不是JSON格式，直接返回原始响应
		respData = map[string]interface{}{
			"raw_response": string(respBody),
		}
	}
	
	// 构建结果
	result := map[string]interface{}{
		"success":         resp.StatusCode >= 200 && resp.StatusCode < 300,
		"status_code":     resp.StatusCode,
		"status":          resp.Status,
		"request_url":     url,
		"request_method":  apiConfig.Method,
		"request_params":  params,
		"request_headers": headers,
		"response":        respData,
		"response_time":   time.Since(startTime),
	}
	
	logrus.Infof("API测试成功，配置ID: %d, URL: %s, 状态码: %d", apiConfig.ID, url, resp.StatusCode)
	
	return result, nil
}

// GetByCode 根据编码获取API配置
func (s *APIConfigService) GetByCode(code string) (*model.APIConfig, error) {
	db := database.GetDB()
	
	var apiConfig model.APIConfig
	if err := db.Where("code = ? AND status = 1", code).First(&apiConfig).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("API配置不存在或已禁用")
		}
		logrus.Errorf("根据编码获取API配置失败: %v", err)
		return nil, err
	}
	
	return &apiConfig, nil
}

// GetByCompanyID 根据公司ID获取API配置列表
func (s *APIConfigService) GetByCompanyID(companyID uint) ([]model.APIConfig, error) {
	db := database.GetDB()
	
	var apiConfigs []model.APIConfig
	if err := db.Where("company_id = ? AND status = 1", companyID).Order("created_at DESC").Find(&apiConfigs).Error; err != nil {
		logrus.Errorf("根据公司ID获取API配置列表失败: %v", err)
		return nil, err
	}
	
	return apiConfigs, nil
}
