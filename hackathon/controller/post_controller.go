package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hackathon/dao/mysql"
	"hackathon/models"
	"hackathon/response"
	"hackathon/service"
	"hackathon/util"
	"net/http"
	"strconv"
)

// CreateTag 创建版块接口
// @Summary 创建版块接口
// @Description 只有身份为管理员才能创建，传name和description 两个参数
// @Tags 版块相关接口
// @Accept application/json
// @Produce application/json
// @Param BearToken header string false "Bearer 用户令牌"
// @Param object body models.ParamCreateTag false "参数"
// @Success 200 {object} response.ResponseData
// @Router /auth/tag [post]
func CreateTag(ctx *gin.Context) {
	//判断权限
	user, _ := ctx.Get("user")
	if user.(models.User).Role != 1 {
		response.Response(ctx, http.StatusOK, response.Error, nil, "不是管理员，权限不足")
		return
	}
	p := new(models.ParamCreateTag)
	if err := ctx.ShouldBindJSON(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}
	code := service.CreateTag(p)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	response.Success(ctx, gin.H{"tag": p}, "创建成功")
}

// RetrieveTags 查询版块接口
// @Summary 查询版块接口
// @Description 查询出所有已创建版块
// @Tags 版块相关接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData
// @Router /tags [get]
func RetrieveTags(ctx *gin.Context) {
	var tags []models.Tag
	mysql.RetrieveArrByStruct(&tags, models.Tag{})
	response.Success(ctx, gin.H{"tags": tags}, "查询成功")
}

// DeleteTag 删除版块接口
// @Summary 删除版块接口
// @Description 只有身份为管理员才能删除，在路径中传id
// @Tags 版块相关接口
// @Accept application/json
// @Produce application/json
// @Param BearToken header string false "Bearer 用户令牌"
// @Param id path int true "版块id"
// @Success 200 {object} response.ResponseData
// @Router /auth/tag/{id} [delete]
func DeleteTag(ctx *gin.Context) {
	//判断权限
	user, _ := ctx.Get("user")
	if user.(models.User).Role != 1 {
		response.Response(ctx, http.StatusOK, response.Error, nil, "不是管理员，权限不足")
		return
	}
	id, _ := strconv.Atoi(ctx.Param("id"))
	var tag models.Tag
	mysql.RetrieveByID(&tag, uint(id))
	err := mysql.Delete(&tag)
	if err != nil {
		response.Response(ctx, http.StatusOK, response.CodeServerBusy, nil, "删除失败,该tag不存在")
		return
	}
	response.Success(ctx, gin.H{"tag": tag}, "删除成功")
}

// RetrievePost 获取帖子接口
// @Summary 获取帖子接口
// @Description 接收路径中的版块id，并随机查询一条该版块下的帖子
// @Tags 树洞相关接口
// @Accept application/json
// @Produce application/json
// @Param id path int true "版块id"
// @Success 200 {object} response.ResponseData
// @Router /post/{id} [get]
func RetrievePost(ctx *gin.Context) {
	tagId, _ := strconv.Atoi(ctx.Param("id"))
	var tag models.Tag
	mysql.RetrieveByID(&tag, uint(tagId))
	if tag.Name == "" {
		response.Response(ctx, http.StatusOK, response.Error, nil, "无效版块")
		return
	}
	// TODO 用redis实现不重复拿数据
	var count models.Post
	mysql.DB.Last(&count)
	fmt.Println(count.ID)
	var post models.Post
	for {
		mysql.DB.Where("tag_id = ? AND id = ? AND status = ?", tagId, util.RandomNumber(1, int(count.ID)), 1).Find(&post)
		if post.Content != "" {
			break
		}
	}
	response.Success(ctx, gin.H{"post": post}, "查询成功")
}

// CreatePost 发贴接口
// @Summary 发贴接口
// @Description 接收tag_id,title,content三个参数，需要token
// @Tags 树洞相关接口
// @Accept application/json
// @Produce application/json
// @Param BearToken header string false "Bearer 用户令牌"
// @Param object body models.ParamCreatePost false "参数"
// @Success 200 {object} response.ResponseData
// @Router /auth/post [post]
func CreatePost(ctx *gin.Context) {
	p := new(models.ParamCreatePost)
	if err := ctx.ShouldBindJSON(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}
	user, _ := ctx.Get("user")
	code := service.CreatePost(p, user.(models.User).ID)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	response.Success(ctx, gin.H{"post": p}, "创建成功")
}

// DeletePost 删贴接口
// @Summary 删贴接口
// @Description 接收路径中的帖子id并删除该帖，只有发贴者和管理员有权操作，需要token
// @Tags 树洞相关接口
// @Accept application/json
// @Produce application/json
// @Param BearToken header string false "Bearer 用户令牌"
// @Param id path int true "版块id"
// @Success 200 {object} response.ResponseData
// @Router /auth/post/{id} [delete]
func DeletePost(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var post models.Post
	mysql.RetrieveByID(&post, uint(id))
	//判断权限
	user, _ := ctx.Get("user")
	if user.(models.User).Role != 1 && user.(models.User).ID != post.UserId {
		response.Response(ctx, http.StatusOK, response.Error, nil, "权限不足")
		return
	}
	err := mysql.Delete(&post)
	if err != nil {
		response.Response(ctx, http.StatusOK, response.CodeServerBusy, nil, "删除失败,该帖不存在")
		return
	}
	response.Success(ctx, gin.H{"post": post}, "删除成功")
}
