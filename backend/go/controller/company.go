package controller

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CompanyController 集团公司控制器
type CompanyController struct {
	companyService *service.CompanyService
}

// NewCompanyController 创建集团公司控制器
func NewCompanyController(companyService *service.CompanyService) *CompanyController {
	return &CompanyController{
		companyService: companyService,
	}
}

// List 获取集团公司列表
func (c *CompanyController) List(ctx *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	// 获取查询参数
	name := ctx.Query("name")
	code := ctx.Query("code")
	typeStr := ctx.Query("type")
	
	var companyType int
	if typeStr != "" {
		companyType, _ = strconv.Atoi(typeStr)
	}
	
	// 调用服务层获取列表
	companies, total, err := c.companyService.List(page, pageSize, name, code, companyType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取集团公司列表失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取集团公司列表成功",
		"data": gin.H{
			"list":      companies,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建集团公司
func (c *CompanyController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var company model.Company
	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层创建集团公司
	if err := c.companyService.Create(&company); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建集团公司失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建集团公司成功",
		"data":    company,
	})
}

// Get 获取集团公司详情
func (c *CompanyController) Get(ctx *gin.Context) {
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
	company, err := c.companyService.Get(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取集团公司详情失败",
			"data":    nil,
		})
		return
	}
	
	if company == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "集团公司不存在",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取集团公司详情成功",
		"data":    company,
	})
}

// Update 更新集团公司
func (c *CompanyController) Update(ctx *gin.Context) {
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
	var company model.Company
	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	company.ID = uint(id)
	
	// 调用服务层更新集团公司
	if err := c.companyService.Update(&company); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新集团公司失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新集团公司成功",
		"data":    company,
	})
}

// Delete 删除集团公司
func (c *CompanyController) Delete(ctx *gin.Context) {
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
	
	// 调用服务层删除集团公司
	if err := c.companyService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除集团公司失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除集团公司成功",
		"data":    nil,
	})
}

// GetTree 获取集团公司树形结构
func (c *CompanyController) GetTree(ctx *gin.Context) {
	// 调用服务层获取树形结构
	tree, err := c.companyService.GetTree()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取集团公司树形结构失败",
			"data":    nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取集团公司树形结构成功",
		"data":    tree,
	})
}
