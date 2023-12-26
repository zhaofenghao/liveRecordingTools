package e

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Ret       int         `json:"ret"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"ret_data"`
	Timestamp int64       `json:"timestamp"`
}

func (g *Gin) Res(httpCode, ret, errCode int, msg string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code:      errCode,
		Msg:       msg,
		Ret:       ret,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

func (g *Gin) PureRes(httpCode, ret, errCode int, msg string, data interface{}) {
	g.C.PureJSON(httpCode, Response{
		Code:      errCode,
		Msg:       msg,
		Ret:       ret,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// 成功
func (g *Gin) Success(msg string, data interface{}) {
	g.Res(http.StatusOK, 0, 1, msg, data)
}

func (g *Gin) SuccessPure(msg string, data interface{}) {
	g.PureRes(http.StatusOK, 0, 1, msg, data)
}

// 可自定义错误码的报错,
func (g *Gin) Fail(code int, msg string, data interface{}) {
	g.Res(http.StatusOK, 1, code, msg, data)
}

// 错误
func (g *Gin) Error(msg string) {
	g.Res(http.StatusOK, 1, -1, msg, nil)
}
