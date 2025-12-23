package database

import (
	"fmt"
	"time"

	"github.com/ddoalistdownload/backend/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// InitMySQL 初始化MySQL连接
// 如果 forceReset 为 true，则会先删除数据库再重新创建
func InitMySQL(cfg *config.MySQLConfig, forceReset bool) error {
	// 先连接到MySQL服务器，不指定数据库名
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
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

	// 连接到MySQL服务器
	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logrus.Errorf("连接MySQL服务器失败: %v", err)
		return err
	}

	// 如果需要重置，先删除数据库
	if forceReset {
		logrus.Warnf("正在执行数据库重置: DROP DATABASE %s", cfg.DBName)
		dropDB := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", cfg.DBName)
		if err := rootDB.Exec(dropDB).Error; err != nil {
			logrus.Errorf("删除数据库失败: %v", err)
			return err
		}
	}

	// 创建数据库
	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s_unicode_ci;",
		cfg.DBName,
		cfg.Charset,
		cfg.Charset,
	)
	if err := rootDB.Exec(createDB).Error; err != nil {
		logrus.Errorf("创建数据库失败: %v", err)
		return err
	}

	// 关闭临时连接
	sqlRootDB, _ := rootDB.DB()
	sqlRootDB.Close()

	// 构建完整DSN（包含数据库名）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	// 连接到指定数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logrus.Errorf("连接数据库失败: %v", err)
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
