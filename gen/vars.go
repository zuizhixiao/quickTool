package gen


var (
	routerTemplate = `
package router

import (
	"github.com/gin-gonic/gin"
	v1 "TEMPLATE/api/v1"
	MIDDLEWARE_IMPORT
)

func InitRouter(Router *gin.RouterGroup) {
ROUTER_TEMP
}
`

	apiTemp = `package v1

import (
	"TEMPLATE/service/SERVICE_NAME"
	"TEMPLATE/types"
	"TEMPLATE/types/response"
	"github.com/gin-gonic/gin"
)
FUNC_LIST
`
	validTemplate = `	var VAR_STRUCT types.STRUCT_E
	err := c.ShouldBind(&VAR_STRUCT)
	if err != nil {
		response.ParamError("参数缺失", c)
		return
	}
`
	functionTemplate = `
func FUNC_NAME(c *gin.Context) {VALID_TEMP
	err IS_DEFINE= SERVICE_NAME.FUNC_NAME(VAR_STRUCT)
	if err != nil {
		response.DbError(err.Error(), c)
		return
	}

	response.Success("", c)
	return
}
`

	serviceTemp = `package SERVICE_NAME

import (
	"TEMPLATE/types"
)
FUNC_LIST
`

	serviceFunctionTemplate = `
func FUNC_NAME(PARAM_TEMP) error {
	//add your code ...
	return nil
}
`

	middlewareTemplate = `package middleware

import (
	"github.com/gin-gonic/gin"
)
FUNC_LIST
`

	middlewareFuncTemplate = `
func FUNC_NAME() gin.HandlerFunc {
	return func(c *gin.Context) {
		//edit your code...
		
		c.Next()
	}
}`

	DockerfileTemplate = `FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /var/www

ADD . .

RUN go build

FROM alpine

ENV TZ Asia/Shanghai

WORKDIR /app/project

COPY --from=builder /var/www ./

EXPOSE 8002

RUN chmod +x /app/project/TEMPLATE

CMD ["./TEMPLATE"]
`
)
