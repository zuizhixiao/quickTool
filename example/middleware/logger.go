package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"TEMPLATE/config"
	"os"
	"time"
)

const leastDay = 7

// CustomFormatter 定义一个自定义的格式化器
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 只返回日志消息内容
	return []byte(fmt.Sprintf("[%s]%s\n", entry.Time.Format("2006-01-02 15:04:05"), entry.Message)), nil
}

func Logger() *logrus.Logger {

	logFilePath := config.GVA_CONFIG.Log.Path
	logSaveDay := config.GVA_CONFIG.Log.Day

	if logSaveDay < leastDay {
		logSaveDay = leastDay
	}

	//创建文件夹
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(logFilePath, os.ModePerm)
		os.Chmod(logFilePath, os.ModePerm)
	}

	logger := logrus.New()

	//设置日志输出
	//logger.Out = src

	//设置日记级别
	logger.SetLevel(logrus.DebugLevel)

	logWriter, _ := rotatelogs.New(
		logFilePath+"/%Y-%m-%d.log",
		// 生成软链，指向最新的日志文件
		//rotatelogs.WithLinkName(fileName),
		// 设置最长保存时间
		rotatelogs.WithMaxAge(leastDay * 24*time.Hour),//7*24*time.Hour
		// 设置日志切割间隔时间
		//rotatelogs.WithRotationTime(24*time.Hour),
	)

	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfhook := lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 新增hook
	logger.AddHook(lfhook)

	return logger
}

func LoggerToFile() gin.HandlerFunc  {
	logger := Logger()
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

func SqlLogger() *logrus.Logger {
	logFilePath := config.GVA_CONFIG.Log.Path
	logSaveDay := config.GVA_CONFIG.Log.Day

	if logSaveDay < leastDay {
		logSaveDay = leastDay
	}

	//创建文件夹
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(logFilePath, os.ModePerm)
		os.Chmod(logFilePath, os.ModePerm)
	}

	logger := logrus.New()

	//设置日记级别
	logger.SetLevel(logrus.InfoLevel)

	logWriter, _ := rotatelogs.New(
		logFilePath+"/%Y-%m-%d.log",
		// 生成软链，指向最新的日志文件
		//rotatelogs.WithLinkName(fileName),
		// 设置最长保存时间
		rotatelogs.WithMaxAge(leastDay * 24*time.Hour),//7*24*time.Hour
		// 设置日志切割间隔时间
		//rotatelogs.WithRotationTime(24*time.Hour),
	)

	logger.SetOutput(logWriter)
	logger.SetFormatter(&CustomFormatter{})
	return logger
}