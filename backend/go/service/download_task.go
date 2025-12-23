package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// DownloadTaskService 下载任务服务
type DownloadTaskService struct{}

// NewDownloadTaskService 创建下载任务服务实例
func NewDownloadTaskService() *DownloadTaskService {
	return &DownloadTaskService{}
}

// List 获取下载任务列表
func (s *DownloadTaskService) List(page, pageSize int, companyID, userID uint, taskName, taskType, status string, currentUser model.User, roleCodes []string) ([]model.DownloadTask, int64, error) {
	db := database.GetDB()

	var downloadTasks []model.DownloadTask
	var total int64

	// 构建查询
	query := db.Model(&model.DownloadTask{})

	// 权限控制：如果不是管理员，只能看自己的任务
	isAdmin := false
	for _, code := range roleCodes {
		if code == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		query = query.Where("user_id = ?", currentUser.ID)
	} else if userID > 0 {
		// 管理员可以指定具体某一用户
		query = query.Where("user_id = ?", userID)
	}

	// 添加其它查询条件
	if companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}
	if taskName != "" {
		query = query.Where("task_name LIKE ?", "%"+taskName+"%")
	}
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		logrus.Errorf("获取下载任务总数失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Company").Preload("User").Preload("APIConfig").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&downloadTasks).Error; err != nil {
		logrus.Errorf("获取下载任务列表失败: %v", err)
		return nil, 0, err
	}

	return downloadTasks, total, nil
}

// Get 获取下载任务详情
func (s *DownloadTaskService) Get(id uint) (*model.DownloadTask, error) {
	db := database.GetDB()

	var downloadTask model.DownloadTask
	if err := db.Preload("Company").Preload("User").Preload("APIConfig").First(&downloadTask, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("下载任务不存在")
		}
		logrus.Errorf("获取下载任务详情失败: %v", err)
		return nil, err
	}

	return &downloadTask, nil
}

// Create 创建下载任务
func (s *DownloadTaskService) Create(downloadTask *model.DownloadTask) error {
	db := database.GetDB()

	// 检查公司是否存在
	var company model.Company
	if err := db.First(&company, downloadTask.CompanyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("公司不存在")
		}
		logrus.Errorf("检查公司是否存在失败: %v", err)
		return err
	}

	// 检查用户是否存在
	var user model.User
	if err := db.First(&user, downloadTask.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		logrus.Errorf("检查用户是否存在失败: %v", err)
		return err
	}

	// 检查API配置是否存在
	var apiConfig model.APIConfig
	if err := db.First(&apiConfig, downloadTask.APIConfigID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("API配置不存在")
		}
		logrus.Errorf("检查API配置是否存在失败: %v", err)
		return err
	}

	// 设置默认值
	if downloadTask.Status == "" {
		downloadTask.Status = "pending"
	}
	if downloadTask.Progress == 0 {
		downloadTask.Progress = 0
	}

	// 创建下载任务
	if err := db.Create(downloadTask).Error; err != nil {
		logrus.Errorf("创建下载任务失败: %v", err)
		return err
	}

	// 异步执行下载任务
	go s.ExecuteTask(downloadTask)

	return nil
}

// Delete 删除下载任务
func (s *DownloadTaskService) Delete(id uint) error {
	db := database.GetDB()

	// 检查下载任务是否存在
	var downloadTask model.DownloadTask
	if err := db.First(&downloadTask, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("下载任务不存在")
		}
		logrus.Errorf("获取下载任务失败: %v", err)
		return err
	}

	// 软删除下载任务
	if err := db.Delete(&downloadTask).Error; err != nil {
		logrus.Errorf("删除下载任务失败: %v", err)
		return err
	}

	// 软删除相关的下载结果
	if err := db.Where("task_id = ?", id).Delete(&model.DownloadResult{}).Error; err != nil {
		logrus.Errorf("删除下载结果失败: %v", err)
		return err
	}

	return nil
}

// ExecuteTask 执行下载任务
func (s *DownloadTaskService) ExecuteTask(task *model.DownloadTask) {
	db := database.GetDB()

	// 更新任务状态为运行中
	task.Status = "running"
	task.Progress = 10
	if err := db.Save(task).Error; err != nil {
		logrus.Errorf("更新任务状态失败: %v", err)
		return
	}

	// 获取API配置
	var apiConfig model.APIConfig
	if err := db.First(&apiConfig, task.APIConfigID).Error; err != nil {
		logrus.Errorf("获取API配置失败: %v", err)
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("获取API配置失败: %v", err)
		if err := db.Save(task).Error; err != nil {
			logrus.Errorf("更新任务状态失败: %v", err)
		}
		return
	}

	// 解析参数
	var params map[string]interface{}
	if task.Params != "" {
		if err := json.Unmarshal([]byte(task.Params), &params); err != nil {
			logrus.Errorf("解析任务参数失败: %v", err)
			task.Status = "failed"
			task.ErrorMsg = fmt.Sprintf("参数格式错误: %v", err)
			if err := db.Save(task).Error; err != nil {
				logrus.Errorf("更新任务状态失败: %v", err)
			}
			return
		}
	}

	// 解析请求头
	var headers map[string]string
	if apiConfig.Headers != "" {
		if err := json.Unmarshal([]byte(apiConfig.Headers), &headers); err != nil {
			logrus.Errorf("解析请求头失败: %v", err)
			task.Status = "failed"
			task.ErrorMsg = fmt.Sprintf("请求头格式错误: %v", err)
			if err := db.Save(task).Error; err != nil {
				logrus.Errorf("更新任务状态失败: %v", err)
			}
			return
		}
	}

	// 构建请求URL
	baseURL := strings.TrimSuffix(apiConfig.BaseURL, "/")
	path := strings.TrimPrefix(apiConfig.Path, "/")
	url := baseURL + "/" + path

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

	if apiConfig.Method == "GET" {
		req, err = http.NewRequest(apiConfig.Method, url, nil)
	} else {
		// 对于非GET请求，将参数转换为JSON
		jsonParams, err := json.Marshal(params)
		if err != nil {
			logrus.Errorf("参数转换为JSON失败: %v", err)
			task.Status = "failed"
			task.ErrorMsg = fmt.Sprintf("参数转换为JSON失败: %v", err)
			if err := db.Save(task).Error; err != nil {
				logrus.Errorf("更新任务状态失败: %v", err)
			}
			return
		}
		req, err = http.NewRequest(apiConfig.Method, url, strings.NewReader(string(jsonParams)))
	}

	if err != nil {
		logrus.Errorf("创建请求失败: %v", err)
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("创建请求失败: %v", err)
		if err := db.Save(task).Error; err != nil {
			logrus.Errorf("更新任务状态失败: %v", err)
		}
		return
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
		logrus.Errorf("发送请求失败: %v", err)
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("发送请求失败: %v", err)
		if err := db.Save(task).Error; err != nil {
			logrus.Errorf("更新任务状态失败: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := s.readResponseBody(resp)
	if err != nil {
		logrus.Errorf("读取响应失败: %v", err)
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("读取响应失败: %v", err)
		if err := db.Save(task).Error; err != nil {
			logrus.Errorf("更新任务状态失败: %v", err)
		}
		return
	}

	// 解析响应
	var respData map[string]interface{}
	if err := json.Unmarshal(respBody, &respData); err != nil {
		// 如果响应不是JSON格式，直接返回原始响应
		respData = map[string]interface{}{
			"raw_response": string(respBody),
		}
	}

	// 更新任务进度
	task.Progress = 80
	if err := db.Save(task).Error; err != nil {
		logrus.Errorf("更新任务进度失败: %v", err)
		return
	}

	// 保存下载结果
	resultJSON, err := json.Marshal(respData)
	if err != nil {
		logrus.Errorf("转换结果为JSON失败: %v", err)
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("转换结果为JSON失败: %v", err)
		if err := db.Save(task).Error; err != nil {
			logrus.Errorf("更新任务状态失败: %v", err)
		}
		return
	}

	// 创建下载结果
	downloadResult := model.DownloadResult{
		TaskID:    task.ID,
		Data:      string(resultJSON),
		DataType:  "json",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&downloadResult).Error; err != nil {
		logrus.Errorf("创建下载结果失败: %v", err)
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("创建下载结果失败: %v", err)
		if err := db.Save(task).Error; err != nil {
			logrus.Errorf("更新任务状态失败: %v", err)
		}
		return
	}

	// 更新任务结果
	task.Status = "success"
	task.Progress = 100
	task.Result = string(resultJSON)
	task.FileName = fmt.Sprintf("download_%d.json", task.ID)
	task.FileSize = int64(len(resultJSON))
	if err := db.Save(task).Error; err != nil {
		logrus.Errorf("更新任务结果失败: %v", err)
		return
	}

	logrus.Infof("下载任务执行成功，任务ID: %d, 文件名: %s", task.ID, task.FileName)
}

// readResponseBody 读取响应体
func (s *DownloadTaskService) readResponseBody(resp *http.Response) ([]byte, error) {
	// 限制响应体大小为100MB
	maxBodySize := int64(100 * 1024 * 1024)

	// 检查响应体大小
	if resp.ContentLength > maxBodySize {
		return nil, errors.New("响应体过大")
	}

	// 读取响应体
	body := make([]byte, 0, resp.ContentLength)
	buf := make([]byte, 4096)
	var totalRead int64

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			totalRead += int64(n)
			if totalRead > maxBodySize {
				return nil, errors.New("响应体过大")
			}
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return body, nil
}

// GetResult 获取下载结果
func (s *DownloadTaskService) GetResult(taskID uint) (*model.DownloadResult, error) {
	db := database.GetDB()

	var result model.DownloadResult
	if err := db.Where("task_id = ?", taskID).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("下载结果不存在")
		}
		logrus.Errorf("获取下载结果失败: %v", err)
		return nil, err
	}

	return &result, nil
}

// GetTaskByUserID 根据用户ID获取下载任务列表
func (s *DownloadTaskService) GetTaskByUserID(userID uint) ([]model.DownloadTask, error) {
	db := database.GetDB()

	var tasks []model.DownloadTask
	if err := db.Where("user_id = ?", userID).Order("created_at DESC").Limit(10).Find(&tasks).Error; err != nil {
		logrus.Errorf("根据用户ID获取下载任务列表失败: %v", err)
		return nil, err
	}

	return tasks, nil
}
