package model

import (
	"time"
)

// DownloadTask 下载任务模型
type DownloadTask struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyID   uint      `gorm:"not null" json:"company_id"`        // 公司ID
	UserID      uint      `gorm:"not null" json:"user_id"`           // 用户ID
	APIConfigID uint      `gorm:"not null" json:"api_config_id"`     // API配置ID
	TaskName    string    `gorm:"size:100;not null" json:"task_name"` // 任务名称
	TaskType    string    `gorm:"size:20;not null" json:"task_type"`  // 任务类型：list, detail
	Params      string    `gorm:"type:text" json:"params"`           // 请求参数（JSON格式）
	Status      string    `gorm:"size:20;default:'pending'" json:"status"` // 任务状态：pending, running, success, failed
	Progress    int       `gorm:"default:0" json:"progress"`         // 任务进度（0-100）
	Result      string    `gorm:"type:text" json:"result"`           // 任务结果（JSON格式）
	FileURL     string    `gorm:"size:255" json:"file_url"`          // 下载文件URL
	FileName    string    `gorm:"size:100" json:"file_name"`         // 下载文件名称
	FileSize    int64     `gorm:"default:0" json:"file_size"`         // 文件大小（字节）
	ErrorMsg    string    `gorm:"type:text" json:"error_msg"`        // 错误信息
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Company     Company   `gorm:"foreignKey:CompanyID" json:"company"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	APIConfig   APIConfig `gorm:"foreignKey:APIConfigID" json:"api_config"`
}

// TableName 设置表名
func (DownloadTask) TableName() string {
	return "download_task"
}

// DownloadResult 下载结果模型
type DownloadResult struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `gorm:"not null" json:"task_id"`           // 下载任务ID
	Data       string    `gorm:"type:longtext" json:"data"`        // 下载数据（JSON格式）
	DataType   string    `gorm:"size:20;not null" json:"data_type"`  // 数据类型：json, excel, csv
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	DownloadTask DownloadTask `gorm:"foreignKey:TaskID" json:"download_task"`
}

// TableName 设置表名
func (DownloadResult) TableName() string {
	return "download_result"
}