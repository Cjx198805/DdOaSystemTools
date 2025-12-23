package service

import (
	"fmt"
	"testing"

	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/ddoalistdownload/backend/util"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to connect database: %v\n", err)
		panic(err)
	}

	// 迁移模型
	db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{}, &model.FieldPermission{}, &model.DataDictionary{})
	database.DB = db
}

func TestCheckFieldEditable(t *testing.T) {
	setupTestDB()
	db := database.GetDB()
	svc := NewFieldPermissionService()

	// 准备测试数据
	userID := uint(1)
	roleID := uint(1)
	module := "user"
	field := "password"

	// 关联用户和角色
	db.Create(&model.UserRole{UserID: userID, RoleID: roleID})

	t.Run("DefaultAllow", func(t *testing.T) {
		// 既没有特殊权限也没有字典限制，默认允许编辑
		editable, err := svc.CheckFieldEditable(userID, module, field)
		if err != nil {
			t.Fatalf("CheckFieldEditable failed: %v", err)
		}
		if !editable {
			t.Error("Expected editable to be true by default")
		}
	})

	t.Run("DictionaryRestrict", func(t *testing.T) {
		// 设置数据字典为不可编辑
		dict := &model.DataDictionary{
			Module:   module,
			Field:    field,
			Editable: util.IntPtr(0),
			Status:   util.IntPtr(1),
		}
		if err := db.Create(dict).Error; err != nil {
			t.Fatalf("Failed to create dictionary: %v", err)
		}

		editable, err := svc.CheckFieldEditable(userID, module, field)
		if err != nil {
			t.Fatalf("CheckFieldEditable failed: %v", err)
		}
		if editable {
			t.Errorf("Expected editable to be false after dictionary restriction. Dict: %+v", dict)
		}
	})

	t.Run("SpecialOverride", func(t *testing.T) {
		// 设置特殊的编辑权限，应当覆盖字典限制
		db.Create(&model.FieldPermission{
			RoleID:      roleID,
			Module:      module,
			Field:       field,
			SpecialEdit: util.IntPtr(1),
		})

		editable, err := svc.CheckFieldEditable(userID, module, field)
		if err != nil {
			t.Fatalf("CheckFieldEditable failed: %v", err)
		}
		if !editable {
			t.Error("Expected editable to be true with special override")
		}
	})
}
