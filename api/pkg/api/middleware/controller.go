package middleware

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// 记录控制器执行时间
func Consuming() gin.HandlerFunc {
	return ControllerConsuming()
}

func ControllerConsuming() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("startTime", time.Now())
		logrus.Info("start controller middleware")
		c.Next()
	}
}
