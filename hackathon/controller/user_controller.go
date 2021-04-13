package controller

import (
	"github.com/gin-gonic/gin"
	"hackathon/dao/mysql"
	"hackathon/models"
	"hackathon/response"
	"hackathon/service"
	"net/http"
	"strconv"
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
	id, _ := strconv.Atoi(ctx.Param("id"))
	var user models.User
	mysql.RetrieveByID(&user, uint(id))
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": models.ToUserDto(user)}})
}

func EditInfo(ctx *gin.Context) {
	p := new(models.ParamEditInfo)
	if err := ctx.ShouldBind(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}
	user, _ := ctx.Get("user")
	code := service.EditInfo(p, user.(models.User).ID)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	response.Success(ctx, gin.H{"data": p}, "修改成功")
}
