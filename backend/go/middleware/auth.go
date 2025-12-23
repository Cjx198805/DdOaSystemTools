package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ddoalistdownload/backend/config"
	"github.com/ddoalistdownload/backend/database"
	"github.com/ddoalistdownload/backend/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret 获取JWT密钥
func getJWTSecret() []byte {
	if config.GlobalConfig != nil && config.GlobalConfig.Server.JWTSecret != "" {
		return []byte(config.GlobalConfig.Server.JWTSecret)
	}
	return []byte("ddoalistdownload-secret-key")
}

// Claims 自定义JWT声明
// 包含用户ID、用户名和角色ID列表
// 用于在JWT令牌中存储用户信息
// 作者: cjx
// 邮箱: xx4125517@126.com
// 时间: 2025-12-22 14:30:00
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleIDs  []uint `json:"role_ids"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
// 参数: userID - 用户ID, username - 用户名, roleIDs - 角色ID列表
// 返回: token字符串和错误信息
// 作者: cjx
// 邮箱: xx4125517@126.com
// 时间: 2025-12-22 14:30:00
func GenerateToken(userID uint, username string, roleIDs []uint) (string, error) {
	// 创建声明
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleIDs:  roleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间设置为72小时
			ExpiresAt: jwt.NewNumericDate(jwt.TimeFunc().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(jwt.TimeFunc()),
			NotBefore: jwt.NewNumericDate(jwt.TimeFunc()),
			Issuer:    "ddoalistdownload",
			Subject:   username,
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT令牌
// 参数: tokenString - JWT令牌字符串
// 返回: Claims对象和错误信息
// 作者: cjx
// 邮箱: xx4125517@126.com
// 时间: 2025-12-22 14:30:00
func ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// AuthMiddleware JWT认证中间件
// 用于验证JWT令牌，确保只有认证用户才能访问API端点
// 作者: cjx
// 邮箱: xx4125517@126.com
// 时间: 2025-12-22 14:30:00
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌格式错误",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 解析JWT令牌
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证令牌",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleIDs", claims.RoleIDs)

		c.Next()
	}
}

// PermissionMiddleware 权限检查中间件
// 用于检查用户是否有权限访问某个API端点
// 参数: requiredRole - 所需的角色代码
// 作者: cjx
// 邮箱: xx4125517@126.com
// 时间: 2025-12-22 14:30:00
func PermissionMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色ID列表
		roleIDs, exists := c.Get("roleIDs")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无法获取用户角色信息",
				"data":    nil,
			})
			c.Abort()
			return
		}

		roleIDList, ok := roleIDs.([]uint)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "用户角色信息格式错误",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 查询角色信息
		db := database.GetDB()
		var roles []model.Role
		if err := db.Where("id IN ?", roleIDList).Find(&roles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询角色信息失败",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查用户是否具有所需角色
		hasPermission := false
		for _, role := range roles {
			if role.Code == requiredRole || role.Code == "admin" { // admin角色拥有所有权限
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "没有权限访问该资源",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
