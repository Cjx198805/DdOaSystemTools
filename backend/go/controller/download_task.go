package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DownloadTaskController 下载任务控制器
type DownloadTaskController struct {
	downloadTaskService *service.DownloadTaskService
}

// NewDownloadTaskController 创建下载任务控制器实例
func NewDownloadTaskController() *DownloadTaskController {
	return &DownloadTaskController{
		downloadTaskService: service.NewDownloadTaskService(),
	}
}

// List 获取下载任务列表
// @Summary 获取下载任务列表
// @Description 分页获取下载任务列表
// @Tags 下载任务管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param company_id query uint false "公司ID"
// @Param user_id query uint false "用户ID"
// @Param task_name query string false "任务名称"
// @Param task_type query string false "任务类型"
// @Param status query string false "任务状态"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/download-task [get]
func (c *DownloadTaskController) List(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	companyIDStr := ctx.Query("company_id")
	userIDStr := ctx.Query("user_id")
	taskName := ctx.Query("task_name")
	taskType := ctx.Query("task_type")
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

	// 调用服务层获取列表
	downloadTasks, total, err := c.downloadTaskService.List(page, pageSize, companyID, userID, taskName, taskType, status)
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
		"message": "获取下载任务列表成功",
		"data": gin.H{
			"list":      downloadTasks,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取下载任务详情
// @Summary 获取下载任务详情
// @Description 根据ID获取下载任务详情
// @Tags 下载任务管理
// @Accept json
// @Produce json
// @Param id path uint true "下载任务ID"
// @Success 200 {object} model.DownloadTask
// @Router /api/v1/download-task/{id} [get]
func (c *DownloadTaskController) Get(ctx *gin.Context) {
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
	downloadTask, err := c.downloadTaskService.Get(uint(id))
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
		"message": "获取下载任务详情成功",
		"data":    downloadTask,
	})
}

// Create 创建下载任务
// @Summary 创建下载任务
// @Description 创建新的下载任务
// @Tags 下载任务管理
// @Accept json
// @Produce json
// @Param download_task body model.DownloadTask true "下载任务信息"
// @Success 200 {object} model.DownloadTask
// @Router /api/v1/download-task [post]
func (c *DownloadTaskController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var downloadTask model.DownloadTask
	if err := ctx.ShouldBindJSON(&downloadTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层创建
	if err := c.downloadTaskService.Create(&downloadTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建下载任务成功",
		"data":    downloadTask,
	})
}

// Delete 删除下载任务
// @Summary 删除下载任务
// @Description 根据ID删除下载任务
// @Tags 下载任务管理
// @Accept json
// @Produce json
// @Param id path uint true "下载任务ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/download-task/{id} [delete]
func (c *DownloadTaskController) Delete(ctx *gin.Context) {
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
	if err := c.downloadTaskService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除下载任务成功",
		"data":    nil,
	})
}

// GetResult 获取下载结果
// @Summary 获取下载结果
// @Description 根据任务ID获取下载结果
// @Tags 下载任务管理
// @Accept json
// @Produce json
// @Param task_id path uint true "下载任务ID"
// @Success 200 {object} model.DownloadResult
// @Router /api/v1/download-task/{task_id}/result [get]
func (c *DownloadTaskController) GetResult(ctx *gin.Context) {
	// 获取ID参数
	taskIDStr := ctx.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层获取结果
	result, err := c.downloadTaskService.GetResult(uint(taskID))
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
		"message": "获取下载结果成功",
		"data":    result,
	})
}

// GetTaskByUserID 根据用户ID获取下载任务列表
// @Summary 根据用户ID获取下载任务列表
// @Description 根据用户ID获取最近的下载任务列表
// @Tags 下载任务管理
// @Accept json
// @Produce json
// @Param user_id path uint true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/download-task/user/{user_id} [get]
func (c *DownloadTaskController) GetTaskByUserID(ctx *gin.Context) {
	// 获取ID参数
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层获取任务列表
	tasks, err := c.downloadTaskService.GetTaskByUserID(uint(userID))
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
		"message": "获取用户下载任务列表成功",
		"data":    tasks,
	})
}