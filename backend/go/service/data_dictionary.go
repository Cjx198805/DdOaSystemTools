package service

import (
	"errors"

	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// DataDictionaryService 数据字典服务
type DataDictionaryService struct{}

// NewDataDictionaryService 创建数据字典服务实例
func NewDataDictionaryService() *DataDictionaryService {
	return &DataDictionaryService{}
}

// List 获取数据字典列表
func (s *DataDictionaryService) List(page, pageSize int, module, field string, status int) ([]model.DataDictionary, int64, error) {
	db := database.GetDB()

	var dataDictionaries []model.DataDictionary
	var total int64

	// 构建查询
	query := db.Model(&model.DataDictionary{})

	// 添加查询条件
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if field != "" {
		query = query.Where("field LIKE ?", "%"+field+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		logrus.Errorf("获取数据字典总数失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&dataDictionaries).Error; err != nil {
		logrus.Errorf("获取数据字典列表失败: %v", err)
		return nil, 0, err
	}

	return dataDictionaries, total, nil
}

// Get 获取数据字典详情
func (s *DataDictionaryService) Get(id uint) (*model.DataDictionary, error) {
	db := database.GetDB()

	var dataDictionary model.DataDictionary
	if err := db.First(&dataDictionary, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("数据字典不存在")
		}
		logrus.Errorf("获取数据字典详情失败: %v", err)
		return nil, err
	}

	return &dataDictionary, nil
}

// Create 创建数据字典
func (s *DataDictionaryService) Create(dataDictionary *model.DataDictionary) error {
	db := database.GetDB()

	// 检查是否已存在相同的模块和字段
	var existing model.DataDictionary
	result := db.Where("module = ? AND field = ?", dataDictionary.Module, dataDictionary.Field).First(&existing)
	if result.Error == nil {
		return errors.New("该模块的该字段数据字典已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查数据字典是否存在失败: %v", result.Error)
		return result.Error
	}

	// 设置默认值
	if dataDictionary.Status == nil {
		dataDictionary.Status = util.IntPtr(1)
	}
	if dataDictionary.DataType == "" {
		dataDictionary.DataType = "string"
	}
	if dataDictionary.Required == nil {
		dataDictionary.Required = util.IntPtr(0)
	}
	if dataDictionary.Editable == nil {
		dataDictionary.Editable = util.IntPtr(1)
	}

	// 创建数据字典
	if err := db.Create(dataDictionary).Error; err != nil {
		logrus.Errorf("创建数据字典失败: %v", err)
		return err
	}

	return nil
}

// Update 更新数据字典
func (s *DataDictionaryService) Update(dataDictionary *model.DataDictionary) error {
	db := database.GetDB()

	// 检查数据字典是否存在
	var existing model.DataDictionary
	if err := db.First(&existing, dataDictionary.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("数据字典不存在")
		}
		logrus.Errorf("获取数据字典失败: %v", err)
		return err
	}

	// 检查是否已存在相同的模块和字段（排除当前记录）
	var duplicate model.DataDictionary
	result := db.Where("module = ? AND field = ? AND id != ?", dataDictionary.Module, dataDictionary.Field, dataDictionary.ID).First(&duplicate)
	if result.Error == nil {
		return errors.New("该模块的该字段数据字典已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logrus.Errorf("检查数据字典是否存在失败: %v", result.Error)
		return result.Error
	}

	// 更新数据字典
	if err := db.Save(dataDictionary).Error; err != nil {
		logrus.Errorf("更新数据字典失败: %v", err)
		return err
	}

	return nil
}

// Delete 删除数据字典
func (s *DataDictionaryService) Delete(id uint) error {
	db := database.GetDB()

	// 检查数据字典是否存在
	var dataDictionary model.DataDictionary
	if err := db.First(&dataDictionary, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("数据字典不存在")
		}
		logrus.Errorf("获取数据字典失败: %v", err)
		return err
	}

	// 软删除数据字典
	if err := db.Delete(&dataDictionary).Error; err != nil {
		logrus.Errorf("删除数据字典失败: %v", err)
		return err
	}

	return nil
}

// GetByModule 根据模块获取数据字典
func (s *DataDictionaryService) GetByModule(module string) ([]model.DataDictionary, error) {
	db := database.GetDB()

	var dataDictionaries []model.DataDictionary
	if err := db.Where("module = ? AND status = 1", module).Find(&dataDictionaries).Error; err != nil {
		logrus.Errorf("根据模块获取数据字典失败: %v", err)
		return nil, err
	}

	return dataDictionaries, nil
}

// GetByModuleAndField 根据模块和字段获取数据字典
func (s *DataDictionaryService) GetByModuleAndField(module, field string) (*model.DataDictionary, error) {
	db := database.GetDB()

	var dataDictionary model.DataDictionary
	result := db.Where("module = ? AND field = ? AND status = 1", module, field).First(&dataDictionary)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("数据字典不存在")
		}
		logrus.Errorf("根据模块和字段获取数据字典失败: %v", result.Error)
		return nil, result.Error
	}

	return &dataDictionary, nil
}
