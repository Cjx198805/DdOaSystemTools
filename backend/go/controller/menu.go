package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// MenuController 菜单控制器
type MenuController struct {
	menuService *service.MenuService
}

// NewMenuController 创建菜单控制器实例
func NewMenuController() *MenuController {
	return &MenuController{
		menuService: service.NewMenuService(),
	}
}

// List 获取菜单列表
// @Summary 获取菜单列表
// @Description 分页获取菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "菜单名称"
// @Param path query string false "菜单路径"
// @Param status query int false "状态 1:启用 0:禁用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/menu [get]
func (c *MenuController) List(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	name := ctx.Query("name")
	path := ctx.Query("path")
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
	menus, total, err := c.menuService.List(page, pageSize, name, path, status)
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
		"message": "获取菜单列表成功",
		"data": gin.H{
			"list":      menus,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取菜单详情
// @Summary 获取菜单详情
// @Description 根据ID获取菜单详情
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path uint true "菜单ID"
// @Success 200 {object} model.Menu
// @Router /api/v1/menu/{id} [get]
func (c *MenuController) Get(ctx *gin.Context) {
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
	menu, err := c.menuService.Get(uint(id))
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
		"message": "获取菜单详情成功",
		"data":    menu,
	})
}

// Create 创建菜单
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu body model.Menu true "菜单信息"
// @Success 200 {object} model.Menu
// @Router /api/v1/menu [post]
func (c *MenuController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var menu model.Menu
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层创建
	if err := c.menuService.Create(&menu); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建菜单成功",
		"data":    menu,
	})
}

// Update 更新菜单
// @Summary 更新菜单
// @Description 更新已有菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path uint true "菜单ID"
// @Param menu body model.Menu true "菜单信息"
// @Success 200 {object} model.Menu
// @Router /api/v1/menu/{id} [put]
func (c *MenuController) Update(ctx *gin.Context) {
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
	var menu model.Menu
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 设置ID
	menu.ID = uint(id)
	
	// 调用服务层更新
	if err := c.menuService.Update(&menu); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新菜单成功",
		"data":    menu,
	})
}

// Delete 删除菜单
// @Summary 删除菜单
// @Description 根据ID删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path uint true "菜单ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/menu/{id} [delete]
func (c *MenuController) Delete(ctx *gin.Context) {
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
	if err := c.menuService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除菜单成功",
		"data":    nil,
	})
}

// GetTree 获取菜单树形结构
// @Summary 获取菜单树形结构
// @Description 获取菜单树形结构
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/menu/tree [get]
func (c *MenuController) GetTree(ctx *gin.Context) {
	// 调用服务层获取树形结构
	tree, err := c.menuService.GetTree()
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
		"message": "获取菜单树形结构成功",
		"data":    tree,
	})
}

// GetByParentID 根据父ID获取子菜单
// @Summary 根据父ID获取子菜单
// @Description 根据父ID获取子菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param parent_id path uint true "父菜单ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/menu/parent/{parent_id} [get]
func (c *MenuController) GetByParentID(ctx *gin.Context) {
	// 获取父ID参数
	parentIDStr := ctx.Param("parent_id")
	parentID, err := strconv.ParseUint(parentIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "父ID参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层获取子菜单
	menus, err := c.menuService.GetByParentID(uint(parentID))
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
		"message": "获取子菜单成功",
		"data":    menus,
	})
}

// GetAll 获取所有菜单
// @Summary 获取所有菜单
// @Description 获取所有菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param status query int false "状态 1:启用 0:禁用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/menu/all [get]
func (c *MenuController) GetAll(ctx *gin.Context) {
	// 获取状态参数
	statusStr := ctx.Query("status")
	
	status := -1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}
	
	// 调用服务层获取所有菜单
	menus, err := c.menuService.GetAll(status)
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
		"message": "获取所有菜单成功",
		"data":    menus,
	})
}
