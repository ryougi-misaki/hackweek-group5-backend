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

func CreateTag(ctx *gin.Context) {
	//判断权限
	user, _ := ctx.Get("user")
	if user.(models.User).Role != 1 {
		response.Response(ctx, http.StatusOK, response.Error, nil, "不是管理员，权限不足")
		return
	}
	p := new(models.ParamCreateTag)
	if err := ctx.ShouldBind(p); err != nil {
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

func RetrieveTags(ctx *gin.Context) {
	var tags []models.Tag
	mysql.RetrieveArrByStruct(&tags, models.Tag{})
	response.Success(ctx, gin.H{"tags": tags}, "查询成功")
}

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

func RetrievePost(ctx *gin.Context) {
	// TODO tag_id不存在？
	var count models.Post
	mysql.DB.Last(&count)
	fmt.Println(count.ID)
	var post models.Post
	for {
		mysql.DB.Where("tag_id = ? AND id = ?", ctx.Param("id"), util.RandomNumber(1, int(count.ID))).Find(&post)
		if post.Content != "" {
			break
		}
	}
	response.Success(ctx, gin.H{"post": post}, "查询成功")
}

func CreatePost(ctx *gin.Context) {
	p := new(models.ParamCreatePost)
	if err := ctx.ShouldBind(p); err != nil {
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
