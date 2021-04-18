package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"hackathon/dao/mysql"
	"hackathon/models"
	"hackathon/response"
	"hackathon/service"
	"net/http"
	"strconv"
)

// Register 注册接口
// @Summary 注册接口
// @Description 目前无短信验证，只需传name，telephone和password三个参数
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamRegister false "comment"
// @Success 200 {object} response.ResponseData
// @Router /register [post]
func Register(ctx *gin.Context) {
	p := new(models.ParamRegister)
	if err := ctx.ShouldBindJSON(p); err != nil {
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

// Login 登入接口
// @Summary 登入接口
// @Description 电话号码+密码登入，返回token,id
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamLogin false "参数"
// @Success 200 {object} response.ResponseData
// @Router /login [post]
func Login(ctx *gin.Context) {
	//获取参数
	p := new(models.ParamLogin)
	if err := ctx.ShouldBindJSON(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}

	//数据验证
	token, id, code := service.Login(p)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token, "id": id}, "登入成功")
}

// Info 查询用户信息接口
// @Summary 查询用户信息接口
// @Description 需要获得一个string id
// @Tags 用户相关接口
// @Accept  json
// @Produce  json
// @Param id path int true "用户id"
// @Success 200 {object} response.ResponseData
// @Router /info/{id} [get]
func Info(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var user models.User
	mysql.RetrieveByID(&user, uint(id))
	response.Success(ctx, gin.H{"user": models.ToUserDto(user)}, "查询成功")
}

// EditInfo 编辑用户信息接口
// @Summary 编辑用户信息接口
// @Description 提交表单，更新个人信息，需要token
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body models.ParamEditInfo false "参数"
// @Success 200 {object} response.ResponseData
// @Router /auth/me [put]
func EditInfo(ctx *gin.Context) {
	p := new(models.ParamEditInfo)
	if err := ctx.ShouldBindJSON(p); err != nil {
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

// ChangePwd 修改用户密码接口
// @Summary 修改用户密码接口
// @Description 提交表单，更新密码，需要token
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body models.ParamChangePwd false "参数"
// @Success 200 {object} response.ResponseData
// @Router /auth/pwd [put]
func ChangePwd(ctx *gin.Context) {
	p := new(models.ParamChangePwd)
	if err := ctx.ShouldBindJSON(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}
	user, _ := ctx.Get("user")
	if err := bcrypt.CompareHashAndPassword([]byte(user.(models.User).Password), []byte(p.Password)); err != nil {
		response.Response(ctx, http.StatusOK, response.CodePwdWrong, nil, response.GetErrMsg(response.CodePwdWrong))
		return
	}
	code := service.ChangePwd(p, user.(models.User).ID)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	response.Success(ctx, gin.H{"data": p}, "修改成功")
}
