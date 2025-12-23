package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RecoverMiddleware 异常捕获中间件
func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印堆栈信息
				logrus.Errorf("Panic recovered: %v\nStack: %s", err, string(debug.Stack()))

				// 返回标准化错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误，请联系管理员",
					"data":    nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
