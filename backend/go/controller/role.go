package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RoleController 角色控制器
type RoleController struct {
	roleService *service.RoleService
}

// NewRoleController 创建角色控制器实例
func NewRoleController() *RoleController {
	return &RoleController{
		roleService: service.NewRoleService(),
	}
}

// List 获取角色列表
// @Summary 获取角色列表
// @Description 分页获取角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "角色名称"
// @Param code query string false "角色编码"
// @Param status query int false "状态 1:启用 0:禁用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/role [get]
func (c *RoleController) List(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
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
	
	status := -1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}
	
	// 调用服务层获取列表
	roles, total, err := c.roleService.List(page, pageSize, name, code, status)
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
		"message": "获取角色列表成功",
		"data": gin.H{
			"list":      roles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path uint true "角色ID"
// @Success 200 {object} model.Role
// @Router /api/v1/role/{id} [get]
func (c *RoleController) Get(ctx *gin.Context) {
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
	role, err := c.roleService.Get(uint(id))
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
		"message": "获取角色详情成功",
		"data":    role,
	})
}

// Create 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role body model.Role true "角色信息"
// @Success 200 {object} model.Role
// @Router /api/v1/role [post]
func (c *RoleController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var role model.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层创建
	if err := c.roleService.Create(&role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建角色成功",
		"data":    role,
	})
}

// Update 更新角色
// @Summary 更新角色
// @Description 更新已有角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path uint true "角色ID"
// @Param role body model.Role true "角色信息"
// @Success 200 {object} model.Role
// @Router /api/v1/role/{id} [put]
func (c *RoleController) Update(ctx *gin.Context) {
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
	var role model.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 设置ID
	role.ID = uint(id)
	
	// 调用服务层更新
	if err := c.roleService.Update(&role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新角色成功",
		"data":    role,
	})
}

// Delete 删除角色
// @Summary 删除角色
// @Description 根据ID删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path uint true "角色ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/role/{id} [delete]
func (c *RoleController) Delete(ctx *gin.Context) {
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
	if err := c.roleService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除角色成功",
		"data":    nil,
	})
}

// GetMenus 获取角色菜单
// @Summary 获取角色菜单
// @Description 根据角色ID获取菜单列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path uint true "角色ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/role/{id}/menus [get]
func (c *RoleController) GetMenus(ctx *gin.Context) {
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
	
	// 调用服务层获取菜单
	menus, err := c.roleService.GetMenus(uint(id))
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
		"message": "获取角色菜单成功",
		"data":    menus,
	})
}

// AssignMenus 分配菜单
// @Summary 分配菜单
// @Description 为角色分配菜单
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path uint true "角色ID"
// @Param menus body struct{MenuIDs []uint `json:"menu_ids" binding:"required"`} true "菜单ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/role/{id}/assign-menus [put]
func (c *RoleController) AssignMenus(ctx *gin.Context) {
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
		MenuIDs []uint `json:"menu_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层分配菜单
	if err := c.roleService.AssignMenus(uint(id), req.MenuIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "分配菜单成功",
		"data":    nil,
	})
}
