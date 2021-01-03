package gin_middleware

import (
	"time"

	golibVar "github.com/feitianlove/golib/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const RequestIDKey = "RequestID"

//web访问日志
func AccessLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		id := GetRequestID(c)
		c.Set(RequestIDKey, id)
		c.Header("X-Client", golibVar.LocalIP)
		c.Header("X-Request-Id", id)
		c.Next()
		client := c.ClientIP()
		latency := time.Since(t)
		status := c.Writer.Status()
		path := c.Request.URL.Path
		contentType := c.ContentType()
		staffName := c.GetString("Staffname")
		log.WithFields(logrus.Fields{
			"request_id":   id,
			"user":         staffName,
			"client":       client,
			"path":         path,
			"content_type": contentType,
			"status":       status,
			"latency":      latency}).Info("access")
	}
}

// 获取request id
func GetRequestID(c *gin.Context) string {
	//尝试请求req提供的request id
	id := c.Request.Header.Get("X-Request-Id")
	if id == "" {
		//尝试获取智能网关的request id
		id = c.Request.Header.Get("X-Rio-Seq")
		if id == "" {
			// 自己生成uuid1
			id = uuid.New().String()
		}
	}
	return id
}
