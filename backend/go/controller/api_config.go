package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// APIConfigController API配置控制器
type APIConfigController struct {
	apiConfigService *service.APIConfigService
}

// NewAPIConfigController 创建API配置控制器实例
func NewAPIConfigController() *APIConfigController {
	return &APIConfigController{
		apiConfigService: service.NewAPIConfigService(),
	}
}

// List 获取API配置列表
// @Summary 获取API配置列表
// @Description 分页获取API配置列表
// @Tags API配置管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param company_id query uint false "公司ID"
// @Param name query string false "配置名称"
// @Param code query string false "配置编码"
// @Param status query int false "状态 1:启用 0:禁用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/api-config [get]
func (c *APIConfigController) List(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	companyIDStr := ctx.Query("company_id")
	name := ctx.Query("name")
	code := ctx.Query("code")
	statusStr := ctx.Query("status")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	
	companyID := uint(0)
	if companyIDStr != "" {
		companyIDUint64, err := strconv.ParseUint(companyIDStr, 10, 32)
		if err == nil {
			companyID = uint(companyIDUint64)
		}
	}
	
	status := -1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}
	
	// 调用服务层获取列表
	apiConfigs, total, err := c.apiConfigService.List(page, pageSize, companyID, name, code, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取API配置列表失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取API配置列表成功",
		"data": gin.H{
			"list":      apiConfigs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取API配置详情
// @Summary 获取API配置详情
// @Description 根据ID获取API配置详情
// @Tags API配置管理
// @Accept json
// @Produce json
// @Param id path uint true "API配置ID"
// @Success 200 {object} model.APIConfig
// @Router /api/v1/api-config/{id} [get]
func (c *APIConfigController) Get(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层获取详情
	apiConfig, err := c.apiConfigService.Get(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	if apiConfig == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "API配置不存在",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取API配置详情成功",
		"data":    apiConfig,
	})
}

// Create 创建API配置
// @Summary 创建API配置
// @Description 创建新的API配置
// @Tags API配置管理
// @Accept json
// @Produce json
// @Param api_config body model.APIConfig true "API配置信息"
// @Success 200 {object} model.APIConfig
// @Router /api/v1/api-config [post]
func (c *APIConfigController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var apiConfig model.APIConfig
	if err := ctx.ShouldBindJSON(&apiConfig); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层创建
	if err := c.apiConfigService.Create(&apiConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建API配置成功",
		"data":    apiConfig,
	})
}

// Update 更新API配置
// @Summary 更新API配置
// @Description 更新已有的API配置
// @Tags API配置管理
// @Accept json
// @Produce json
// @Param id path uint true "API配置ID"
// @Param api_config body model.APIConfig true "API配置信息"
// @Success 200 {object} model.APIConfig
// @Router /api/v1/api-config/{id} [put]
func (c *APIConfigController) Update(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID参数错误",
			"data":    nil,
		})
		return
	}
	
	// 绑定请求参数
	var apiConfig model.APIConfig
	if err := ctx.ShouldBindJSON(&apiConfig); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 设置ID
	apiConfig.ID = uint(id)
	
	// 调用服务层更新
	if err := c.apiConfigService.Update(&apiConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新API配置成功",
		"data":    apiConfig,
	})
}

// Delete 删除API配置
// @Summary 删除API配置
// @Description 根据ID删除API配置
// @Tags API配置管理
// @Accept json
// @Produce json
// @Param id path uint true "API配置ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/api-config/{id} [delete]
func (c *APIConfigController) Delete(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层删除
	if err := c.apiConfigService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除API配置成功",
		"data":    nil,
	})
}

// Test 测试API配置
// @Summary 测试API配置
// @Description 测试API配置是否可用
// @Tags API配置管理
// @Accept json
// @Produce json
// @Param api_config body model.APIConfig true "API配置信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/api-config/test [post]
func (c *APIConfigController) Test(ctx *gin.Context) {
	// 绑定请求参数
	var apiConfig model.APIConfig
	if err := ctx.ShouldBindJSON(&apiConfig); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层测试
	result, err := c.apiConfigService.Test(&apiConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试API配置成功",
		"data":    result,
	})
}
