package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"library/config"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 将读取的body写回，因为body只能读取一次
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 使用自定义的ResponseWriter来捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 获取请求路径和方法
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()

		// 获取用户信息
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")

		// 构造日志字段
		fields := logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"method":       reqMethod,
			"uri":          reqUri,
		}

		if userID != nil {
			fields["user_id"] = userID
			fields["username"] = username
		}

		// 添加请求体（如果不是文件上传）
		if c.Request.Header.Get("Content-Type") != "multipart/form-data" && len(requestBody) > 0 {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, requestBody, "", "  "); err == nil {
				fields["request_body"] = prettyJSON.String()
			}
		}

		// 添加响应体（如果不是文件下载）
		if !bytes.Contains(blw.body.Bytes(), []byte("application/octet-stream")) && blw.body.Len() > 0 {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, blw.body.Bytes(), "", "  "); err == nil {
				fields["response_body"] = prettyJSON.String()
			}
		}

		// 根据状态码记录日志级别
		if statusCode >= 500 {
			logger.WithFields(fields).Error("Server Error")
		} else if statusCode >= 400 {
			logger.WithFields(fields).Warn("Client Error")
		} else {
			logger.WithFields(fields).Info("Request Completed")
		}
	}
}

// bodyLogWriter 自定义ResponseWriter，用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// InitLogger 初始化日志
func InitLogger() *logrus.Logger {
	logger := logrus.New()

	// 设置日志格式为JSON
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置输出
	if config.GlobalConfig.Server.LogPath != "" {
		// TODO: 实现日志文件轮转
		// file, err := os.OpenFile(config.GlobalConfig.Server.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		// if err == nil {
		// 	logger.Out = file
		// }
	}

	// 设置日志级别
	if config.GlobalConfig.Server.Mode == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
