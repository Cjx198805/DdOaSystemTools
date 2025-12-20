package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/utils"
)

// AccessTokenService AccessToken服务
type AccessTokenService struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewAccessTokenService 创建AccessToken服务实例
func NewAccessTokenService() *AccessTokenService {
	return &AccessTokenService{
		DB:    database.GetDB(),
		Redis: database.GetRedis(),
	}
}

// GetAccessToken 获取AccessToken
func (s *AccessTokenService) GetAccessToken(c *gin.Context, companyID uint) (*model.AccessToken, error) {
	// 先从数据库查询
	var accessToken model.AccessToken
	result := s.DB.Where("company_id = ? AND status = 1", companyID).First(&accessToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("未找到有效的AccessToken配置")
		}
		return nil, result.Error
	}

	// 检查是否过期，过期则刷新
	if time.Now().After(accessToken.ExpiresAt) {
		return s.RefreshAccessToken(c, companyID)
	}

	return &accessToken, nil
}

// RefreshAccessToken 刷新AccessToken
func (s *AccessTokenService) RefreshAccessToken(c *gin.Context, companyID uint) (*model.AccessToken, error) {
	// 获取AccessToken配置
	var accessToken model.AccessToken
	result := s.DB.Where("company_id = ? AND status = 1", companyID).First(&accessToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("未找到有效的AccessToken配置")
		}
		return nil, result.Error
	}

	// 调用钉钉API刷新AccessToken
	newToken, expiresIn, err := s.getAccessTokenFromDingTalk(accessToken.AppKey, accessToken.AppSecret)
	if err != nil {
		return nil, err
	}

	// 更新AccessToken信息
	now := time.Now()
	accessToken.AccessToken = newToken
	accessToken.ExpiresIn = expiresIn
	accessToken.ExpiresAt = now.Add(time.Duration(expiresIn) * time.Second)
	accessToken.RefreshAt = now

	// 保存到数据库
	result = s.DB.Save(&accessToken)
	if result.Error != nil {
		return nil, result.Error
	}

	// 保存到Redis
	redisKey := fmt.Sprintf("access_token:%d", companyID)
	s.Redis.Set(context.Background(), redisKey, newToken, time.Duration(expiresIn)*time.Second)

	return &accessToken, nil
}

// CreateAccessToken 创建AccessToken配置
func (s *AccessTokenService) CreateAccessToken(c *gin.Context, req *model.AccessToken) (*model.AccessToken, error) {
	// 检查公司是否存在
	var company model.Company
	if err := s.DB.First(&company, req.CompanyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("公司不存在")
		}
		return nil, err
	}

	// 检查是否已存在AccessToken配置
	var existingToken model.AccessToken
	result := s.DB.Where("company_id = ?", req.CompanyID).First(&existingToken)
	if result.Error == nil {
		return nil, errors.New("该公司已存在AccessToken配置")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 初始获取AccessToken
	accessToken, expiresIn, err := s.getAccessTokenFromDingTalk(req.AppKey, req.AppSecret)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	req.AccessToken = accessToken
	req.ExpiresIn = expiresIn
	req.ExpiresAt = now.Add(time.Duration(expiresIn) * time.Second)
	req.RefreshAt = now
	req.Status = 1

	// 保存到数据库
	result = s.DB.Create(req)
	if result.Error != nil {
		return nil, result.Error
	}

	// 保存到Redis
	redisKey := fmt.Sprintf("access_token:%d", req.CompanyID)
	s.Redis.Set(context.Background(), redisKey, accessToken, time.Duration(expiresIn)*time.Second)

	return req, nil
}

// UpdateAccessToken 更新AccessToken配置
func (s *AccessTokenService) UpdateAccessToken(c *gin.Context, id uint, req *model.AccessToken) (*model.AccessToken, error) {
	// 检查是否存在
	var accessToken model.AccessToken
	result := s.DB.First(&accessToken, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("AccessToken配置不存在")
		}
		return nil, result.Error
	}

	// 更新字段
	accessToken.AppKey = req.AppKey
	accessToken.AppSecret = req.AppSecret
	accessToken.Status = req.Status

	// 如果状态为有效，重新获取AccessToken
	if req.Status == 1 {
		newToken, expiresIn, err := s.getAccessTokenFromDingTalk(req.AppKey, req.AppSecret)
		if err != nil {
			return nil, err
		}

		now := time.Now()
		accessToken.AccessToken = newToken
		accessToken.ExpiresIn = expiresIn
		accessToken.ExpiresAt = now.Add(time.Duration(expiresIn) * time.Second)
		accessToken.RefreshAt = now
	}

	// 保存到数据库
	result = s.DB.Save(&accessToken)
	if result.Error != nil {
		return nil, result.Error
	}

	// 更新或删除Redis
	redisKey := fmt.Sprintf("access_token:%d", accessToken.CompanyID)
	if accessToken.Status == 1 {
		s.Redis.Set(context.Background(), redisKey, accessToken.AccessToken, time.Until(accessToken.ExpiresAt))
	} else {
		s.Redis.Del(context.Background(), redisKey)
	}

	return &accessToken, nil
}

// DeleteAccessToken 删除AccessToken配置
func (s *AccessTokenService) DeleteAccessToken(c *gin.Context, id uint) error {
	// 检查是否存在
	var accessToken model.AccessToken
	result := s.DB.First(&accessToken, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("AccessToken配置不存在")
		}
		return result.Error
	}

	// 删除Redis中的AccessToken
	redisKey := fmt.Sprintf("access_token:%d", accessToken.CompanyID)
	s.Redis.Del(context.Background(), redisKey)

	// 软删除
	result = s.DB.Delete(&accessToken)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetAccessTokenList 获取AccessToken列表
func (s *AccessTokenService) GetAccessTokenList(c *gin.Context, page, pageSize int) ([]model.AccessToken, int64, error) {
	var accessTokens []model.AccessToken
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询总数
	result := s.DB.Model(&model.AccessToken{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// 查询列表
	result = s.DB.Preload("Company").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&accessTokens)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return accessTokens, total, nil
}

// TestAccessToken 测试AccessToken有效性
func (s *AccessTokenService) TestAccessToken(c *gin.Context, companyID uint) error {
	// 获取AccessToken
	accessToken, err := s.GetAccessToken(c, companyID)
	if err != nil {
		return err
	}

	// 调用钉钉API测试
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/v2/user/get?access_token=%s", accessToken.AccessToken)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("AccessToken无效，HTTP状态码：%d", resp.StatusCode))
	}

	return nil
}

// getAccessTokenFromDingTalk 从钉钉API获取AccessToken
func (s *AccessTokenService) getAccessTokenFromDingTalk(appKey, appSecret string) (string, int, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", appKey, appSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	// 解析响应
	var result map[string]interface{}
	if err := utils.ParseJSON(resp.Body, &result); err != nil {
		return "", 0, err
	}

	// 检查错误码
	errcode, ok := result["errcode"].(float64)
	if !ok || errcode != 0 {
		errmsg, _ := result["errmsg"].(string)
		return "", 0, errors.New(fmt.Sprintf("获取AccessToken失败：%s", errmsg))
	}

	// 提取AccessToken和过期时间
	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", 0, errors.New("获取AccessToken失败：响应格式错误")
	}

	expiresIn, ok := result["expires_in"].(float64)
	if !ok {
		expiresIn = 7200 // 默认2小时
	}

	return accessToken, int(expiresIn), nil
}
