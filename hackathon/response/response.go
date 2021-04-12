package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CodePhoneLength = 1000 + iota
	CodeParamError
	CodePwdLength
	CodePhoneExist
	CodeEncryptError
	CodeUserNotExist
	CodePwdWrong
	CodeServerBusy

	OK    = 0
	Error = 1
)

var codeMsg = map[int]string{
	OK:               "OK",
	Error:            "FAIL",
	CodeParamError:   "各种奇奇怪怪的参数错误",
	CodePhoneLength:  "手机号必须为十一位",
	CodePwdLength:    "密码不得少于6位",
	CodePhoneExist:   "用户已经存在",
	CodeEncryptError: "加密错误",
	CodeUserNotExist: "用户不存在",
	CodePwdWrong:     "密码错误",
	CodeServerBusy:   "服务繁忙",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}

func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": msg})
}

func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, OK, data, msg)
}

func Fail(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, Error, data, msg)
}
