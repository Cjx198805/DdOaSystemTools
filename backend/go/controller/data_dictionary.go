package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DataDictionaryController 数据字典控制器
type DataDictionaryController struct {
	dataDictionaryService *service.DataDictionaryService
}

// NewDataDictionaryController 创建数据字典控制器实例
func NewDataDictionaryController() *DataDictionaryController {
	return &DataDictionaryController{
		dataDictionaryService: service.NewDataDictionaryService(),
	}
}

// List 获取数据字典列表
// @Summary 获取数据字典列表
// @Description 分页获取数据字典列表
// @Tags 数据字典管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param module query string false "模块名称"
// @Param field query string false "字段名称"
// @Param status query int false "状态 1:启用 0:禁用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-dictionary [get]
func (c *DataDictionaryController) List(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	module := ctx.Query("module")
	field := ctx.Query("field")
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
	dataDictionaries, total, err := c.dataDictionaryService.List(page, pageSize, module, field, status)
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
		"message": "获取数据字典列表成功",
		"data": gin.H{
			"list":      dataDictionaries,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取数据字典详情
// @Summary 获取数据字典详情
// @Description 根据ID获取数据字典详情
// @Tags 数据字典管理
// @Accept json
// @Produce json
// @Param id path uint true "数据字典ID"
// @Success 200 {object} model.DataDictionary
// @Router /api/v1/data-dictionary/{id} [get]
func (c *DataDictionaryController) Get(ctx *gin.Context) {
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
	dataDictionary, err := c.dataDictionaryService.Get(uint(id))
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
		"message": "获取数据字典详情成功",
		"data":    dataDictionary,
	})
}

// Create 创建数据字典
// @Summary 创建数据字典
// @Description 创建新的数据字典
// @Tags 数据字典管理
// @Accept json
// @Produce json
// @Param data_dictionary body model.DataDictionary true "数据字典信息"
// @Success 200 {object} model.DataDictionary
// @Router /api/v1/data-dictionary [post]
func (c *DataDictionaryController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var dataDictionary model.DataDictionary
	if err := ctx.ShouldBindJSON(&dataDictionary); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 调用服务层创建
	if err := c.dataDictionaryService.Create(&dataDictionary); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建数据字典成功",
		"data":    dataDictionary,
	})
}

// Update 更新数据字典
// @Summary 更新数据字典
// @Description 更新已有数据字典
// @Tags 数据字典管理
// @Accept json
// @Produce json
// @Param id path uint true "数据字典ID"
// @Param data_dictionary body model.DataDictionary true "数据字典信息"
// @Success 200 {object} model.DataDictionary
// @Router /api/v1/data-dictionary/{id} [put]
func (c *DataDictionaryController) Update(ctx *gin.Context) {
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
	var dataDictionary model.DataDictionary
	if err := ctx.ShouldBindJSON(&dataDictionary); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 设置ID
	dataDictionary.ID = uint(id)

	// 调用服务层更新
	if err := c.dataDictionaryService.Update(&dataDictionary); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新数据字典成功",
		"data":    dataDictionary,
	})
}

// Delete 删除数据字典
// @Summary 删除数据字典
// @Description 根据ID删除数据字典
// @Tags 数据字典管理
// @Accept json
// @Produce json
// @Param id path uint true "数据字典ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-dictionary/{id} [delete]
func (c *DataDictionaryController) Delete(ctx *gin.Context) {
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
	if err := c.dataDictionaryService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除数据字典成功",
		"data":    nil,
	})
}

// GetByModule 根据模块获取数据字典
// @Summary 根据模块获取数据字典
// @Description 根据模块名称获取数据字典列表
// @Tags 数据字典管理
// @Accept json
// @Produce json
// @Param module path string true "模块名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-dictionary/module/{module} [get]
func (c *DataDictionaryController) GetByModule(ctx *gin.Context) {
	// 获取参数
	module := ctx.Param("module")

	// 调用服务层获取数据字典
	dataDictionaries, err := c.dataDictionaryService.GetByModule(module)
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
		"message": "获取数据字典成功",
		"data":    dataDictionaries,
	})
}

// GetByModuleAndField 根据模块和字段获取数据字典
// @Summary 根据模块和字段获取数据字典
// @Description 根据模块名称和字段名称获取数据字典
// @Tags 数据字典管理
// @Accept json
// @Produce json
// @Param module path string true "模块名称"
// @Param field path string true "字段名称"
// @Success 200 {object} model.DataDictionary
// @Router /api/v1/data-dictionary/module/{module}/field/{field} [get]
func (c *DataDictionaryController) GetByModuleAndField(ctx *gin.Context) {
	// 获取参数
	module := ctx.Param("module")
	field := ctx.Param("field")

	// 调用服务层获取数据字典
	dataDictionary, err := c.dataDictionaryService.GetByModuleAndField(module, field)
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
		"message": "获取数据字典成功",
		"data":    dataDictionary,
	})
}