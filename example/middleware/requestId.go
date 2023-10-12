package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestId() gin.HandlerFunc  {
	return func(c *gin.Context) {
		headerName := "X-Request-Id"
		requestUid := c.Request.Header.Get(headerName)
		if requestUid == "" {
			requestUid = genUUID()
		}
		c.Set(headerName, requestUid)
		c.Writer.Header().Set(headerName, requestUid)
		c.Next()
	}
}

func genUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
