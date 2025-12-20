package database

import (
	"fmt"
	"github.com/ddoalistdownload/backend/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	DB *gorm.DB
)

// InitMySQL 初始化MySQL连接
func InitMySQL(cfg *config.MySQLConfig) error {
	// 构建DSN（Data Source Name）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	// 配置GORM日志
	newLogger := logger.New(
		logrus.New(),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 彩色打印
		},
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logrus.Errorf("连接MySQL失败: %v", err)
		return err
	}

	// 获取底层SQL连接池
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("获取SQL连接池失败: %v", err)
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	DB = db
	logrus.Info("MySQL连接初始化成功")
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}
