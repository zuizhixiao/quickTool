package middleware

import (
	"TEMPLATE/config"
	"TEMPLATE/types/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xinliangnote/go-util/mail"
	"runtime/debug"
	"strings"
	"time"
)

func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				DebugStack := ""
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + "<br>"
				}
				subject := "【重要错误】项目出错了！"
				body := strings.ReplaceAll(MailTemplate, "{ErrorMsg}", fmt.Sprintf("%s", err))
				body  = strings.ReplaceAll(body, "{RequestTime}", time.Now().Format("2006-01-02 15:04:05"))
				body  = strings.ReplaceAll(body, "{RequestURL}", c.Request.Method + "  " + c.Request.Host + c.Request.RequestURI)
				body  = strings.ReplaceAll(body, "{RequestUA}", c.Request.UserAgent())
				body  = strings.ReplaceAll(body, "{RequestIP}", c.ClientIP())
				body  = strings.ReplaceAll(body, "{DebugStack}", DebugStack)

				options := &mail.Options{
					MailHost : config.GVA_CONFIG.Email.Host,
					MailPort : config.GVA_CONFIG.Email.Port,
					MailUser : config.GVA_CONFIG.Email.User,
					MailPass : config.GVA_CONFIG.Email.Pass,
					MailTo   : config.GVA_CONFIG.Email.AdminUser,
					Subject  : subject,
					Body     : body,
				}
				_ = mail.Send(options)
				response.SystemError("系统异常", c)
			}
		}()
		c.Next()
	}
}
