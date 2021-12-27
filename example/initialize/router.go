package initialize

import (
	"TEMPLATE/config"
	"TEMPLATE/middleware"
	"TEMPLATE/router"
	"TEMPLATE/types/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RunServer() {
	gin.SetMode(config.GVA_CONFIG.System.Mode)
	//加载路由
	router := Routers()

	address := fmt.Sprintf(":%d", config.GVA_CONFIG.System.Addr)

	//http
	s := &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func Routers() *gin.Engine {
	Router := gin.Default()

	//处理跨域
	Router.Use(cors())

	//日记记录
	Router.Use(middleware.RequestId(), middleware.LoggerToFile()) //, middleware.Exception()

	ApiGroup := Router.Group("")
	router.InitRouter(ApiGroup) //注册用户相关接口路由

	//处理404
	Router.NoMethod(HandleNotFind)
	Router.NoRoute(HandleNotFind)

	return Router
}

//处理404
func HandleNotFind(c *gin.Context)  {
	response.Result(4, "api不存在", "", c)
}

//跨域
func cors() gin.HandlerFunc  {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-ca-stage")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//过滤options请求
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		//处理请求
		c.Next()
	}
}
