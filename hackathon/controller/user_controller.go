package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hackathon/models"
	"hackathon/response"
	"hackathon/service"
	"net/http"
)

func Register(ctx *gin.Context) {
	p := new(models.ParamRegister)
	if err := ctx.ShouldBind(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}
	//数据验证
	code := service.Register(p)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	response.Success(ctx, nil, "注册成功")
}

func Login(ctx *gin.Context) {
	//获取参数
	p := new(models.ParamLogin)
	if err := ctx.ShouldBind(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}

	//数据验证
	token, code := service.Login(p)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token}, "登入成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	//ctx.JSON(http.StatusOK,gin.H{"code":200,"data":gin.H{"user":dto.ToUserDto(user.(model.User))}})
	fmt.Println(user)
}

func SendCode(ctx *gin.Context) {
	telephone := ctx.PostForm("telephone")
	if len(telephone) != 11 {
		response.Fail(ctx, nil, response.GetErrMsg(response.CodePhoneLength))
		return
	}
	ok := service.SendSmsCode(telephone)
	if !ok {
		response.Response(ctx, http.StatusOK, response.Error, nil, "发送验证码失败")
		return
	}
	response.Success(ctx, nil, "发送成功")
}

func ChangePwd(ctx *gin.Context) {
	p := new(models.ParamChangePwd)
	if err := ctx.ShouldBind(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}
	//数据验证
	code := service.ChangePwd(p)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	response.Success(ctx, nil, "注册成功")
}
