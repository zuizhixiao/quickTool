package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zjswh/go-tool/utils"
)

func RequestId() gin.HandlerFunc  {
	return func(c *gin.Context) {
		headerName := "X-Request-Id"
		requestUid := c.Request.Header.Get(headerName)
		if requestUid == "" {
			requestUid = utils.GenUUID()
		}
		c.Set(headerName, requestUid)
		c.Writer.Header().Set(headerName, requestUid)
		c.Next()
	}
}
