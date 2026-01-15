package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	MySQL    MySQLConfig
	Redis    RedisConfig
	DingTalk DingTalkConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port      string
	JWTSecret string
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	Charset  string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	AppKey    string
	AppSecret string
}

var GlobalConfig *Config

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 从环境变量或配置文件加载配置
	// 这里简单实现，实际项目中可以使用viper等配置库
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			Username: getEnv("MYSQL_USERNAME", "root"),
			Password: getEnv("MYSQL_PASSWORD", "root"),
			DBName:   getEnv("MYSQL_DBNAME", "dd_oa_download"),
			Charset:  getEnv("MYSQL_CHARSET", "utf8mb4"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		DingTalk: DingTalkConfig{
			AppKey:    getEnv("DINGTALK_APPKEY", ""),
			AppSecret: getEnv("DINGTALK_APPSECRET", ""),
		},
	}
	config.Server.JWTSecret = getEnv("JWT_SECRET", "ddoalistdownload-secret-key")

	GlobalConfig = config
	logrus.Info("配置加载完成")
	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
