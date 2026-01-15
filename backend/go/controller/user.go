package controller

import (
	"net/http"
	"strconv"

	"github.com/ddoalistdownload/backend/middleware"
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// GetCurrentUserInfo 获取当前登录用户信息
func (c *UserController) GetCurrentUserInfo(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
			"data":    nil,
		})
		return
	}

	// 调用服务层获取详情
	user, err := c.userService.Get(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 转换为 Vben Admin 期望的格式
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取用户信息成功",
		"data": gin.H{
			"userId":   user.ID,
			"username": user.Username,
			"realName": user.Nickname,
			"avatar":   "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
			"desc":     user.Email,
			"roles":    []string{"admin"}, // 暂时硬编码
		},
	})
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param company_id query uint false "公司ID"
// @Param username query string false "用户名"
// @Param nickname query string false "昵称"
// @Param status query int false "状态 1:启用 0:禁用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user [get]
func (c *UserController) List(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	companyIDStr := ctx.Query("company_id")
	username := ctx.Query("username")
	nickname := ctx.Query("nickname")
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
	users, total, err := c.userService.List(page, pageSize, companyID, username, nickname, status)
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
		"message": "获取用户列表成功",
		"data": gin.H{
			"list":      users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取用户详情
// @Summary 获取用户详情
// @Description 根据ID获取用户详情
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} model.User
// @Router /api/v1/user/{id} [get]
func (c *UserController) Get(ctx *gin.Context) {
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
	user, err := c.userService.Get(uint(id))
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
		"message": "获取用户详情成功",
		"data":    user,
	})
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body model.User true "用户信息"
// @Success 200 {object} model.User
// @Router /api/v1/user [post]
func (c *UserController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层创建
	if err := c.userService.Create(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 清空密码
	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建用户成功",
		"data":    user,
	})
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新已有用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Param user body model.User true "用户信息"
// @Success 200 {object} model.User
// @Router /api/v1/user/{id} [put]
func (c *UserController) Update(ctx *gin.Context) {
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
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 设置ID
	user.ID = uint(id)

	// 调用服务层更新
	if err := c.userService.Update(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 清空密码
	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新用户成功",
		"data":    user,
	})
}

// Delete 删除用户
// @Summary 删除用户
// @Description 根据ID删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
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
	if err := c.userService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除用户成功",
		"data":    nil,
	})
}

// ResetPassword 重置用户密码
// @Summary 重置用户密码
// @Description 根据ID重置用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Param password body struct{Password string} true "新密码"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/{id}/reset-password [put]
func (c *UserController) ResetPassword(ctx *gin.Context) {
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
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层重置密码
	if err := c.userService.ResetPassword(uint(id), req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "重置密码成功",
		"data":    nil,
	})
}

// UpdatePassword 更新用户密码
// @Summary 更新用户密码
// @Description 更新当前用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param password body struct{OldPassword string `json:"old_password" binding:"required"`;NewPassword string `json:"new_password" binding:"required"`} true "密码信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/update-password [put]
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	// 绑定请求参数
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 从上下文获取用户ID（实际项目中应该从JWT令牌或会话中获取）
	// 这里简化处理，假设用户ID为1
	userID := uint(1)

	// 调用服务层更新密码
	if err := c.userService.UpdatePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新密码成功",
		"data":    nil,
	})
}

// GetRoles 获取用户角色
// @Summary 获取用户角色
// @Description 根据用户ID获取角色列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/{id}/roles [get]
func (c *UserController) GetRoles(ctx *gin.Context) {
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

	// 调用服务层获取角色
	roles, err := c.userService.GetRoles(uint(id))
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
		"message": "获取用户角色成功",
		"data":    roles,
	})
}

// AssignRoles 分配角色
// @Summary 分配角色
// @Description 为用户分配角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Param roles body struct{RoleIDs []uint `json:"role_ids" binding:"required"`} true "角色ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/{id}/assign-roles [put]
func (c *UserController) AssignRoles(ctx *gin.Context) {
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
	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层分配角色
	if err := c.userService.AssignRoles(uint(id), req.RoleIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "分配角色成功",
		"data":    nil,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param login body struct{Username string `json:"username" binding:"required"`;Password string `json:"password" binding:"required"`} true "登录信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	// 绑定请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层登录
	user, roleIDs, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user.ID, user.Username, roleIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成令牌失败",
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}
