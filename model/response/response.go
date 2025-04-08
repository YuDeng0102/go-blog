package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	Error   = 7
	Success = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}
func Ok(c *gin.Context) {
	Result(Success, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(Success, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(Success, data, "success", c)
}

func OkWithDetailed(data interface{}, msg string, c *gin.Context) {
	Result(Success, data, msg, c)
}

func Fail(c *gin.Context) {
	Result(Error, map[string]interface{}{}, "failure", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(Error, map[string]interface{}{}, message, c)
}

func FailWithDetiled(data interface{}, msg string, c *gin.Context) {
	Result(Error, data, msg, c)
}

func NoAuth(message string, c *gin.Context) {
	Result(Error, gin.H{"reload": true}, message, c)
}

func Forbidden(message string, c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code: Error,
		Data: nil,
		Msg:  message,
	})
}
