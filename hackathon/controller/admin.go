package controller

import (
	"github.com/gin-gonic/gin"
	"hackathon/dao/mysql"
	"hackathon/models"
	"hackathon/response"
	"net/http"
	"strconv"
)

// ShowAllPosts 查询帖子接口
// @Summary 查询帖子接口
// @Description 查询出未审核的帖子,只有管理员有权限调，需要token
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} response.ResponseData
// @Router /admin/posts [get]
func ShowAllPosts(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	if user.(models.User).Role != 1 {
		response.Response(ctx, http.StatusOK, response.Error, nil, "不是管理员，权限不足")
		return
	}
	posts := []models.Post{}
	mysql.GetDB().Where("status = ?", 0).Find(&posts)
	response.Success(ctx, gin.H{"posts": posts}, "查询成功")
}

// UpdatePostStatus 审核帖子接口
// @Summary 审核帖子接口
// @Description 审核帖子，路径传id和status；status为1：审核通过，为2：审核不通过。只有管理员有权限调，需要token
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path int true "帖子id"
// @Param status path int true "状态"
// @Success 200 {object} response.ResponseData
// @Router /admin/post/{id}/{status} [put]
func UpdatePostStatus(ctx *gin.Context) {
	//判断权限
	user, _ := ctx.Get("user")
	if user.(models.User).Role != 1 {
		response.Response(ctx, http.StatusOK, response.Error, nil, "不是管理员，权限不足")
		return
	}
	id, _ := strconv.Atoi(ctx.Param("id"))
	status, _ := strconv.Atoi(ctx.Param("status"))
	var post models.Post
	mysql.RetrieveByID(&post, uint(id))
	updateData := &models.Post{
		Status: status,
	}
	err := mysql.Update(post, updateData)
	if err != nil {
		response.Response(ctx, http.StatusOK, response.CodeServerBusy, nil, response.GetErrMsg(response.CodeServerBusy))
		return
	}
	response.Success(ctx, gin.H{"post": post}, "操作成功")
}
