package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FieldPermissionController 字段权限控制器
type FieldPermissionController struct {
	fieldPermissionService *service.FieldPermissionService
}

// NewFieldPermissionController 创建字段权限控制器实例
func NewFieldPermissionController() *FieldPermissionController {
	return &FieldPermissionController{
		fieldPermissionService: service.NewFieldPermissionService(),
	}
}

// List 获取字段权限列表
// @Summary 获取字段权限列表
// @Description 分页获取字段权限列表
// @Tags 字段权限管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param role_id query uint false "角色ID"
// @Param module query string false "模块名称"
// @Param field query string false "字段名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/field-permission [get]
func (c *FieldPermissionController) List(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	roleIDStr := ctx.Query("role_id")
	module := ctx.Query("module")
	field := ctx.Query("field")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	roleID := uint(0)
	if roleIDStr != "" {
		roleIDUint64, err := strconv.ParseUint(roleIDStr, 10, 32)
		if err == nil {
			roleID = uint(roleIDUint64)
		}
	}

	// 调用服务层获取列表
	fieldPermissions, total, err := c.fieldPermissionService.List(page, pageSize, roleID, module, field)
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
		"message": "获取字段权限列表成功",
		"data": gin.H{
			"list":      fieldPermissions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取字段权限详情
// @Summary 获取字段权限详情
// @Description 根据ID获取字段权限详情
// @Tags 字段权限管理
// @Accept json
// @Produce json
// @Param id path uint true "字段权限ID"
// @Success 200 {object} model.FieldPermission
// @Router /api/v1/field-permission/{id} [get]
func (c *FieldPermissionController) Get(ctx *gin.Context) {
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
	fieldPermission, err := c.fieldPermissionService.Get(uint(id))
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
		"message": "获取字段权限详情成功",
		"data":    fieldPermission,
	})
}

// Create 创建字段权限
// @Summary 创建字段权限
// @Description 创建新的字段权限
// @Tags 字段权限管理
// @Accept json
// @Produce json
// @Param field_permission body model.FieldPermission true "字段权限信息"
// @Success 200 {object} model.FieldPermission
// @Router /api/v1/field-permission [post]
func (c *FieldPermissionController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var fieldPermission model.FieldPermission
	if err := ctx.ShouldBindJSON(&fieldPermission); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层创建
	if err := c.fieldPermissionService.Create(&fieldPermission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建字段权限成功",
		"data":    fieldPermission,
	})
}

// Update 更新字段权限
// @Summary 更新字段权限
// @Description 更新已有字段权限
// @Tags 字段权限管理
// @Accept json
// @Produce json
// @Param id path uint true "字段权限ID"
// @Param field_permission body model.FieldPermission true "字段权限信息"
// @Success 200 {object} model.FieldPermission
// @Router /api/v1/field-permission/{id} [put]
func (c *FieldPermissionController) Update(ctx *gin.Context) {
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
	var fieldPermission model.FieldPermission
	if err := ctx.ShouldBindJSON(&fieldPermission); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 设置ID
	fieldPermission.ID = uint(id)

	// 调用服务层更新
	if err := c.fieldPermissionService.Update(&fieldPermission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新字段权限成功",
		"data":    fieldPermission,
	})
}

// Delete 删除字段权限
// @Summary 删除字段权限
// @Description 根据ID删除字段权限
// @Tags 字段权限管理
// @Accept json
// @Produce json
// @Param id path uint true "字段权限ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/field-permission/{id} [delete]
func (c *FieldPermissionController) Delete(ctx *gin.Context) {
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
	if err := c.fieldPermissionService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除字段权限成功",
		"data":    nil,
	})
}

// GetByRoleAndModule 根据角色和模块获取字段权限
// @Summary 根据角色和模块获取字段权限
// @Description 根据角色ID和模块获取字段权限列表
// @Tags 字段权限管理
// @Accept json
// @Produce json
// @Param role_id path uint true "角色ID"
// @Param module path string true "模块名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/field-permission/role/{role_id}/module/{module} [get]
func (c *FieldPermissionController) GetByRoleAndModule(ctx *gin.Context) {
	// 获取参数
	roleIDStr := ctx.Param("role_id")
	module := ctx.Param("module")

	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "角色ID参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层获取字段权限
	fieldPermissions, err := c.fieldPermissionService.GetByRoleAndModule(uint(roleID), module)
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
		"message": "获取字段权限成功",
		"data":    fieldPermissions,
	})
}