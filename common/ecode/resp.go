package ecode

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Response struct {
	RetCode ECode       `json:"ret"`
	ErrMsg  string      `json:"errmsg"`
	Data    interface{} `json:"data,omitempty"`
	Count   int32       `json:"count,omitempty"`
}

func RespOkMsg(c *gin.Context, errMsg string) {
	result := Response{RetCode: OK, ErrMsg: fmt.Sprintf("%s:%s", OK, errMsg)}
	c.JSON(200, result)
	c.Next()
}
func RespOkData(c *gin.Context, data interface{}, errMsg string) {
	result := Response{RetCode: OK, Data: data, ErrMsg: fmt.Sprintf("%s:%s", OK, errMsg)}
	c.JSON(200, result)
	c.Next()
}
func RespOkDataCount(c *gin.Context, data interface{}, count int32, errMsg string) {
	result := Response{RetCode: OK, Data: data, Count: count, ErrMsg: fmt.Sprintf("%s:%s", OK, errMsg)}
	c.JSON(200, result)
	c.Next()
}

func RespErrCode(c *gin.Context, e ECode, errMsg string) {
	result := Response{RetCode: ParamEmpty, ErrMsg: fmt.Sprintf("%s:%s", e, errMsg)}
	c.JSON(200, result)
	c.Next()
}

func RespStatus403(c *gin.Context, e ECode, errMsg string) {
	result := Response{RetCode: ParamEmpty, ErrMsg: fmt.Sprintf("%s:%s", e, errMsg)}
	c.JSON(403, result)
	c.Abort()
}
