package model

import (
	"time"
)

// APITestCase API测试用例模型
type APITestCase struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyID   uint      `gorm:"not null" json:"company_id"`        // 公司ID
	APIConfigID uint      `gorm:"not null" json:"api_config_id"`     // API配置ID
	Name        string    `gorm:"size:100;not null" json:"name"`     // 测试用例名称
	Description string    `gorm:"type:text" json:"description"`     // 测试用例描述
	Params      string    `gorm:"type:text" json:"params"`           // 请求参数（JSON格式）
	Headers     string    `gorm:"type:text" json:"headers"`          // 请求头（JSON格式）
	ExpectedResult string  `gorm:"type:text" json:"expected_result"`  // 期望结果
	Status      int       `gorm:"default:1" json:"status"`             // 状态 1:启用 0:禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Company     Company   `gorm:"foreignKey:CompanyID" json:"company"`
	APIConfig   APIConfig `gorm:"foreignKey:APIConfigID" json:"api_config"`
}

// TableName 设置表名
func (APITestCase) TableName() string {
	return "api_test_case"
}

// APITestHistory API测试历史记录模型
type APITestHistory struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CompanyID       uint      `gorm:"not null" json:"company_id"`        // 公司ID
	UserID          uint      `gorm:"not null" json:"user_id"`           // 用户ID
	APIConfigID     uint      `gorm:"not null" json:"api_config_id"`     // API配置ID
	TestCaseID      uint      `json:"test_case_id"`                        // 测试用例ID（可选）
	Name            string    `gorm:"size:100;not null" json:"name"`     // 测试名称
	Params          string    `gorm:"type:text" json:"params"`           // 请求参数（JSON格式）
	Headers         string    `gorm:"type:text" json:"headers"`          // 请求头（JSON格式）
	ActualResult    string    `gorm:"type:text" json:"actual_result"`    // 实际结果
	Status          string    `gorm:"size:20" json:"status"`             // 测试状态：success, failed
	ResponseTime    int64     `json:"response_time"`                     // 响应时间（毫秒）
	StatusCode      int       `json:"status_code"`                     // HTTP状态码
	ErrorMessage    string    `gorm:"type:text" json:"error_message"`    // 错误信息
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Company         Company   `gorm:"foreignKey:CompanyID" json:"company"`
	User            User      `gorm:"foreignKey:UserID" json:"user"`
	APIConfig       APIConfig `gorm:"foreignKey:APIConfigID" json:"api_config"`
	TestCase        APITestCase `gorm:"foreignKey:TestCaseID" json:"test_case"`
}

// TableName 设置表名
func (APITestHistory) TableName() string {
	return "api_test_history"
}