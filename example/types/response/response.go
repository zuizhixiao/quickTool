package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code         int         `json:"code"`
	Data         interface{} `json:"data"`
	ErrorMessage string      `json:"errorMessage"`
	ErrorCode    int         `json:"errorCode"`
}

func Success(data interface{}, c *gin.Context) {
	Result(0, data, "", c)
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	result := Response{
		200,
		data,
		msg,
		code,
	}
	c.JSON(http.StatusOK, result)
}

func SystemError(msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusInternalServerError, Response{
		http.StatusInternalServerError,
		"",
		msg,
		0,
	})
	c.Abort()
}

func ParamError(message string, c *gin.Context) {
	Result(2, "", message, c)
}

func DbError(message string, c *gin.Context) {
	Result(3, "", message, c)
}

