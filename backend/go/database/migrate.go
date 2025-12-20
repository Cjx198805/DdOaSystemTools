package database

import (
	"github.com/ddoalistdownload/backend/model"
	"github.com/sirupsen/logrus"
)

// MigrateDB 数据库迁移
func MigrateDB() error {
	logrus.Info("开始数据库迁移...")

	// 使用 AutoMigrate 自动创建表
	err := DB.AutoMigrate(
		&model.Company{},
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.Menu{},
		&model.RoleMenu{},
		&model.SSOConfig{},
		&model.AccessToken{},
		&model.APIConfig{},
		&model.Log{},
		&model.FieldPermission{},
		&model.DataDictionary{},
		&model.DownloadTask{},
		&model.DownloadResult{},
		&model.APITestCase{},
		&model.APITestHistory{},
	)

	if err != nil {
		logrus.Errorf("数据库迁移失败: %v", err)
		return err
	}

	logrus.Info("数据库迁移完成")
	return nil
}
