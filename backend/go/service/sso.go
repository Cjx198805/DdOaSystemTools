package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
)

// SSOService 身份验证（免登）服务
type SSOService struct {}

// NewSSOService 创建身份验证（免登）服务
func NewSSOService() *SSOService {
	return &SSOService{}
}

// GetConfig 获取身份验证（免登）配置
func (s *SSOService) GetConfig(companyID uint) (*model.SSOConfig, error) {
	db := database.GetDB()
	
	var config model.SSOConfig
	if err := db.Where("company_id = ?", companyID).First(&config).Error; err != nil {
		logrus.Errorf("获取身份验证（免登）配置失败: %v", err)
		return nil, err
	}
	
	return &config, nil
}

// UpdateConfig 更新身份验证（免登）配置
func (s *SSOService) UpdateConfig(config *model.SSOConfig) error {
	db := database.GetDB()
	
	// 检查是否存在
	var existing model.SSOConfig
	result := db.Where("company_id = ?", config.CompanyID).First(&existing)
	if result.Error != nil {
		// 不存在则创建
		if err := db.Create(config).Error; err != nil {
			logrus.Errorf("创建身份验证（免登）配置失败: %v", err)
			return err
		}
		return nil
	}
	
	// 存在则更新
	config.ID = existing.ID
	if err := db.Save(config).Error; err != nil {
		logrus.Errorf("更新身份验证（免登）配置失败: %v", err)
		return err
	}
	
	return nil
}

// TestSSO 测试身份验证（免登）
func (s *SSOService) TestSSO(companyID uint, code string) (map[string]interface{}, error) {
	db := database.GetDB()
	
	// 获取配置
	var config model.SSOConfig
	if err := db.Where("company_id = ?", companyID).First(&config).Error; err != nil {
		logrus.Errorf("获取身份验证（免登）配置失败: %v", err)
		return nil, err
	}
	
	// 1. 获取access_token
	accessTokenURL := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", config.AppKey, config.AppSecret)
	resp, err := http.Get(accessTokenURL)
	if err != nil {
		logrus.Errorf("获取access_token失败: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	// 解析access_token响应
	var accessTokenResp struct {
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&accessTokenResp); err != nil {
		logrus.Errorf("解析access_token响应失败: %v", err)
		return nil, err
	}
	
	if accessTokenResp.Errcode != 0 {
		logrus.Errorf("获取access_token失败: %s", accessTokenResp.Errmsg)
		return nil, errors.New(fmt.Sprintf("获取access_token失败: %s", accessTokenResp.Errmsg))
	}
	
	// 2. 使用code获取用户信息
	userInfoURL := fmt.Sprintf("https://oapi.dingtalk.com/user/getuserinfo?access_token=%s&code=%s", accessTokenResp.AccessToken, code)
	resp, err = http.Get(userInfoURL)
	if err != nil {
		logrus.Errorf("获取用户信息失败: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	// 解析用户信息响应
	var userInfoResp struct {
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		Userid      string `json:"userid"`
		Name        string `json:"name"`
		Department  []int  `json:"department"`
		Unionid     string `json:"unionid"`
		Openid      string `json:"openid"`
		Avatar      string `json:"avatar"`
		Mobile      string `json:"mobile"`
		Email       string `json:"email"`
		Jobnumber   string `json:"jobnumber"`
		Position    string `json:"position"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfoResp); err != nil {
		logrus.Errorf("解析用户信息响应失败: %v", err)
		return nil, err
	}
	
	if userInfoResp.Errcode != 0 {
		logrus.Errorf("获取用户信息失败: %s", userInfoResp.Errmsg)
		return nil, errors.New(fmt.Sprintf("获取用户信息失败: %s", userInfoResp.Errmsg))
	}
	
	// 返回结果
	result := map[string]interface{}{
		"success": true,
		"message": "免登测试成功",
		"data": map[string]interface{}{
			"company_id": companyID,
			"app_key":    config.AppKey,
			"code":       code,
			"user_info": map[string]interface{}{
				"userid":     userInfoResp.Userid,
				"name":       userInfoResp.Name,
				"department": userInfoResp.Department,
				"unionid":    userInfoResp.Unionid,
				"openid":     userInfoResp.Openid,
				"avatar":     userInfoResp.Avatar,
				"mobile":     userInfoResp.Mobile,
				"email":      userInfoResp.Email,
				"jobnumber":  userInfoResp.Jobnumber,
				"position":   userInfoResp.Position,
			},
		},
	}
	
	logrus.Infof("免登测试成功，公司ID: %d, AppKey: %s, Code: %s, UserID: %s", companyID, config.AppKey, code, userInfoResp.Userid)
	
	return result, nil
}
