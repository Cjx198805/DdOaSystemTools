package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
)

// AccessTokenController AccessToken控制器
type AccessTokenController struct {
	accessTokenService *service.AccessTokenService
}

// NewAccessTokenController 创建AccessToken控制器实例
func NewAccessTokenController() *AccessTokenController {
	return &AccessTokenController{
		accessTokenService: service.NewAccessTokenService(),
	}
}

// GetAccessToken 获取AccessToken
// @Summary 获取AccessToken
// @Description 根据公司ID获取AccessToken
// @Tags AccessToken管理
// @Accept json
// @Produce json
// @Param company_id query uint true "公司ID"
// @Success 200 {object} model.AccessToken
// @Router /api/access-token [get]
func (c *AccessTokenController) GetAccessToken(ctx *gin.Context) {
	// 获取公司ID
	companyIDStr := ctx.Query("company_id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的公司ID"})
		return
	}

	// 调用服务
	accessToken, err := c.accessTokenService.GetAccessToken(ctx, uint(companyID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}

// CreateAccessToken 创建AccessToken配置
// @Summary 创建AccessToken配置
// @Description 创建新的AccessToken配置
// @Tags AccessToken管理
// @Accept json
// @Produce json
// @Param access_token body model.AccessToken true "AccessToken配置"
// @Success 200 {object} model.AccessToken
// @Router /api/access-token [post]
func (c *AccessTokenController) CreateAccessToken(ctx *gin.Context) {
	// 绑定请求
	var req model.AccessToken
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	accessToken, err := c.accessTokenService.CreateAccessToken(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}

// UpdateAccessToken 更新AccessToken配置
// @Summary 更新AccessToken配置
// @Description 更新已有的AccessToken配置
// @Tags AccessToken管理
// @Accept json
// @Produce json
// @Param id path uint true "AccessToken ID"
// @Param access_token body model.AccessToken true "AccessToken配置"
// @Success 200 {object} model.AccessToken
// @Router /api/access-token/{id} [put]
func (c *AccessTokenController) UpdateAccessToken(ctx *gin.Context) {
	// 获取ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 绑定请求
	var req model.AccessToken
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	accessToken, err := c.accessTokenService.UpdateAccessToken(ctx, uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}

// DeleteAccessToken 删除AccessToken配置
// @Summary 删除AccessToken配置
// @Description 删除已有的AccessToken配置
// @Tags AccessToken管理
// @Accept json
// @Produce json
// @Param id path uint true "AccessToken ID"
// @Success 200 {object} map[string]string
// @Router /api/access-token/{id} [delete]
func (c *AccessTokenController) DeleteAccessToken(ctx *gin.Context) {
	// 获取ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 调用服务
	err = c.accessTokenService.DeleteAccessToken(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetAccessTokenList 获取AccessToken列表
// @Summary 获取AccessToken列表
// @Description 分页获取AccessToken列表
// @Tags AccessToken管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/access-token/list [get]
func (c *AccessTokenController) GetAccessTokenList(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 调用服务
	accessTokens, total, err := c.accessTokenService.GetAccessTokenList(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"list":  accessTokens,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// RefreshAccessToken 刷新AccessToken
// @Summary 刷新AccessToken
// @Description 手动刷新AccessToken
// @Tags AccessToken管理
// @Accept json
// @Produce json
// @Param company_id query uint true "公司ID"
// @Success 200 {object} model.AccessToken
// @Router /api/access-token/refresh [post]
func (c *AccessTokenController) RefreshAccessToken(ctx *gin.Context) {
	// 获取公司ID
	companyIDStr := ctx.Query("company_id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的公司ID"})
		return
	}

	// 调用服务
	accessToken, err := c.accessTokenService.RefreshAccessToken(ctx, uint(companyID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}

// TestAccessToken 测试AccessToken有效性
// @Summary 测试AccessToken有效性
// @Description 测试AccessToken是否有效
// @Tags AccessToken管理
// @Accept json
// @Produce json
// @Param company_id query uint true "公司ID"
// @Success 200 {object} map[string]string
// @Router /api/access-token/test [post]
func (c *AccessTokenController) TestAccessToken(ctx *gin.Context) {
	// 获取公司ID
	companyIDStr := ctx.Query("company_id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的公司ID"})
		return
	}

	// 调用服务
	err = c.accessTokenService.TestAccessToken(ctx, uint(companyID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "AccessToken有效"})
}
