package middleware

import (
	"TEMPLATE/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func LoggerToFile() gin.HandlerFunc  {
	logFilePath := config.GVA_CONFIG.Log.Path
	logFileName := config.GVA_CONFIG.Log.Name
	fileName := path.Join(logFilePath, logFileName)

	//创建文件夹
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(logFilePath, os.ModePerm)
		os.Chmod(logFilePath, os.ModePerm)
	}

	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY , os.ModeAppend)
	if err != nil {
		fmt.Println("err:", err.Error())
	}

	logger := logrus.New()

	//设置日志输出
	logger.Out = src

	//设置日记级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})



	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method

		reqUri := c.Request.RequestURI

		status := c.Writer.Status()

		clientIp := c.ClientIP()

		logger.WithFields(logrus.Fields{
			"status" : status ,
			"latencyTime" : latencyTime ,
			"clientIp" : clientIp ,
			"reqMethod" : reqMethod ,
			"reqUri" : reqUri ,
		}).Info()

	}
}
