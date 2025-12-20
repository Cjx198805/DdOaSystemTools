package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// APITestController API测试控制器
type APITestController struct {
	apiTestService *service.APITestService
}

// NewAPITestController 创建API测试控制器实例
func NewAPITestController() *APITestController {
	return &APITestController{
		apiTestService: service.NewAPITestService(),
	}
}

// ListTestCases 获取API测试用例列表
// @Summary 获取API测试用例列表
// @Description 分页获取API测试用例列表
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param company_id query uint false "公司ID"
// @Param api_config_id query uint false "API配置ID"
// @Param name query string false "测试用例名称"
// @Param status query int false "状态 1:启用 0:禁用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/api-test/case [get]
func (c *APITestController) ListTestCases(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	companyIDStr := ctx.Query("company_id")
	apiConfigIDStr := ctx.Query("api_config_id")
	name := ctx.Query("name")
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

	apiConfigID := uint(0)
	if apiConfigIDStr != "" {
		apiConfigIDUint64, err := strconv.ParseUint(apiConfigIDStr, 10, 32)
		if err == nil {
			apiConfigID = uint(apiConfigIDUint64)
		}
	}

	status := -1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}

	// 调用服务层获取列表
	testCases, total, err := c.apiTestService.ListTestCases(page, pageSize, companyID, apiConfigID, name, status)
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
		"message": "获取API测试用例列表成功",
		"data": gin.H{
			"list":      testCases,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetTestCase 获取API测试用例详情
// @Summary 获取API测试用例详情
// @Description 根据ID获取API测试用例详情
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param id path uint true "测试用例ID"
// @Success 200 {object} model.APITestCase
// @Router /api/v1/api-test/case/{id} [get]
func (c *APITestController) GetTestCase(ctx *gin.Context) {
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
	testCase, err := c.apiTestService.GetTestCase(uint(id))
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
		"message": "获取API测试用例详情成功",
		"data":    testCase,
	})
}

// CreateTestCase 创建API测试用例
// @Summary 创建API测试用例
// @Description 创建新的API测试用例
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param test_case body model.APITestCase true "API测试用例信息"
// @Success 200 {object} model.APITestCase
// @Router /api/v1/api-test/case [post]
func (c *APITestController) CreateTestCase(ctx *gin.Context) {
	// 绑定请求参数
	var testCase model.APITestCase
	if err := ctx.ShouldBindJSON(&testCase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层创建
	if err := c.apiTestService.CreateTestCase(&testCase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建API测试用例成功",
		"data":    testCase,
	})
}

// UpdateTestCase 更新API测试用例
// @Summary 更新API测试用例
// @Description 更新已有API测试用例
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param id path uint true "测试用例ID"
// @Param test_case body model.APITestCase true "API测试用例信息"
// @Success 200 {object} model.APITestCase
// @Router /api/v1/api-test/case/{id} [put]
func (c *APITestController) UpdateTestCase(ctx *gin.Context) {
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
	var testCase model.APITestCase
	if err := ctx.ShouldBindJSON(&testCase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 设置ID
	testCase.ID = uint(id)

	// 调用服务层更新
	if err := c.apiTestService.UpdateTestCase(&testCase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新API测试用例成功",
		"data":    testCase,
	})
}

// DeleteTestCase 删除API测试用例
// @Summary 删除API测试用例
// @Description 根据ID删除API测试用例
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param id path uint true "测试用例ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/api-test/case/{id} [delete]
func (c *APITestController) DeleteTestCase(ctx *gin.Context) {
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
	if err := c.apiTestService.DeleteTestCase(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除API测试用例成功",
		"data":    nil,
	})
}

// RunTestCase 执行API测试用例
// @Summary 执行API测试用例
// @Description 根据测试用例ID执行API测试
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param id path uint true "测试用例ID"
// @Param user_id body struct{UserID uint `json:"user_id"`} true "用户ID"
// @Success 200 {object} model.APITestHistory
// @Router /api/v1/api-test/case/{id}/run [post]
func (c *APITestController) RunTestCase(ctx *gin.Context) {
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
	var req struct{
		UserID uint `json:"user_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 获取测试用例
	testCase, err := c.apiTestService.GetTestCase(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用服务层执行测试用例
	testHistory, err := c.apiTestService.RunTestCase(req.UserID, testCase)
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
		"message": "执行API测试用例成功",
		"data":    testHistory,
	})
}

// ListTestHistory 获取API测试历史记录列表
// @Summary 获取API测试历史记录列表
// @Description 分页获取API测试历史记录列表
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param company_id query uint false "公司ID"
// @Param user_id query uint false "用户ID"
// @Param api_config_id query uint false "API配置ID"
// @Param test_case_id query uint false "测试用例ID"
// @Param name query string false "测试名称"
// @Param status query string false "测试状态"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/api-test/history [get]
func (c *APITestController) ListTestHistory(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	companyIDStr := ctx.Query("company_id")
	userIDStr := ctx.Query("user_id")
	apiConfigIDStr := ctx.Query("api_config_id")
	testCaseIDStr := ctx.Query("test_case_id")
	name := ctx.Query("name")
	status := ctx.Query("status")

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

	userID := uint(0)
	if userIDStr != "" {
		userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			userID = uint(userIDUint64)
		}
	}

	apiConfigID := uint(0)
	if apiConfigIDStr != "" {
		apiConfigIDUint64, err := strconv.ParseUint(apiConfigIDStr, 10, 32)
		if err == nil {
			apiConfigID = uint(apiConfigIDUint64)
		}
	}

	testCaseID := uint(0)
	if testCaseIDStr != "" {
		testCaseIDUint64, err := strconv.ParseUint(testCaseIDStr, 10, 32)
		if err == nil {
			testCaseID = uint(testCaseIDUint64)
		}
	}

	// 调用服务层获取列表
	testHistories, total, err := c.apiTestService.ListTestHistory(page, pageSize, companyID, userID, apiConfigID, testCaseID, name, status)
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
		"message": "获取API测试历史记录列表成功",
		"data": gin.H{
			"list":      testHistories,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetTestHistory 获取API测试历史记录详情
// @Summary 获取API测试历史记录详情
// @Description 根据ID获取API测试历史记录详情
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param id path uint true "测试历史记录ID"
// @Success 200 {object} model.APITestHistory
// @Router /api/v1/api-test/history/{id} [get]
func (c *APITestController) GetTestHistory(ctx *gin.Context) {
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
	testHistory, err := c.apiTestService.GetTestHistory(uint(id))
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
		"message": "获取API测试历史记录详情成功",
		"data":    testHistory,
	})
}

// DeleteTestHistory 删除API测试历史记录
// @Summary 删除API测试历史记录
// @Description 根据ID删除API测试历史记录
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param id path uint true "测试历史记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/api-test/history/{id} [delete]
func (c *APITestController) DeleteTestHistory(ctx *gin.Context) {
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
	if err := c.apiTestService.DeleteTestHistory(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除API测试历史记录成功",
		"data":    nil,
	})
}

// ClearTestHistory 清空API测试历史记录
// @Summary 清空API测试历史记录
// @Description 清空指定公司的所有API测试历史记录
// @Tags API测试管理
// @Accept json
// @Produce json
// @Param company_id body struct{CompanyID uint `json:"company_id"`} true "公司ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/api-test/history/clear [post]
func (c *APITestController) ClearTestHistory(ctx *gin.Context) {
	// 绑定请求参数
	var req struct{
		CompanyID uint `json:"company_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层清空测试历史记录
	if err := c.apiTestService.ClearTestHistory(req.CompanyID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "清空API测试历史记录成功",
		"data":    nil,
	})
}