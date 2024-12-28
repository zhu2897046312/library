package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				logger.WithFields(logrus.Fields{
					"error": err,
					"stack": string(debug.Stack()),
				}).Error("Server Error")

				// 如果是开发环境，返回详细错误信息
				var message interface{} = "服务器内部错误"
				if gin.Mode() == gin.DebugMode {
					message = err
				}

				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  message,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ValidationMiddleware 参数验证错误处理中间件
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理验证错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypeBind:
					// 参数绑定错误
					c.JSON(http.StatusBadRequest, gin.H{
						"code": 400,
						"msg":  e.Err.Error(),
					})
				default:
					// 其他错误
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  e.Err.Error(),
					})
				}
				return
			}
		}
	}
}

// ResponseMiddleware 统一响应处理中间件
func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果已经有响应了就不处理
		if c.Writer.Written() {
			return
		}

		// 获取设置的响应数据
		if data, exists := c.Get("response"); exists {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "success",
				"data": data,
			})
			return
		}
	}
}
