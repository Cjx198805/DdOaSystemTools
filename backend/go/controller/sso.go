package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SSOController 身份验证（免登）控制器
type SSOController struct {
	ssoService *service.SSOService
}

// NewSSOController 创建身份验证（免登）控制器
func NewSSOController() *SSOController {
	return &SSOController{
		ssoService: service.NewSSOService(),
	}
}

// GetConfig 获取身份验证（免登）配置
func (c *SSOController) GetConfig(ctx *gin.Context) {
	// 获取公司ID参数
	companyIDStr := ctx.Query("company_id")
	if companyIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "公司ID不能为空",
			"data":    nil,
		})
		return
	}
	
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "公司ID参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层获取配置
	config, err := c.ssoService.GetConfig(uint(companyID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取身份验证（免登）配置失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取身份验证（免登）配置成功",
		"data":    config,
	})
}

// UpdateConfig 更新身份验证（免登）配置
func (c *SSOController) UpdateConfig(ctx *gin.Context) {
	// 绑定请求参数
	var config model.SSOConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层更新配置
	if err := c.ssoService.UpdateConfig(&config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新身份验证（免登）配置失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新身份验证（免登）配置成功",
		"data":    config,
	})
}

// TestSSO 测试身份验证（免登）
func (c *SSOController) TestSSO(ctx *gin.Context) {
	// 绑定请求参数
	var req struct {
		CompanyID uint   `json:"company_id" binding:"required"`
		Code      string `json:"code" binding:"required"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层测试免登
	result, err := c.ssoService.TestSSO(req.CompanyID, req.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "测试身份验证（免登）失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试身份验证（免登）成功",
		"data":    result,
	})
}
