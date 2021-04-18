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

// ShowMyChats 显示聊天列表接口
// @Summary 显示聊天列表接口
// @Description 展示个人聊天列表，需要token
// @Tags 聊天列表相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} response.ResponseData
// @Router /auth/chats [get]
func ShowMyChats(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	chats := []models.ChatRecord{}
	mysql.GetDB().Where("user_id = ?", user.(models.User).ID).Find(&chats)
	response.Success(ctx, gin.H{"chats": chats}, "查询成功")
}

// AddChatRecord 更新聊天列表接口
// @Summary 更新聊天列表接口
// @Description 增加聊天列表，如果存在，则更新最后一条聊天记录，需要token
// @Tags 聊天列表相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body models.ParamAddChatRecord false "参数"
// @Success 200 {object} response.ResponseData
// @Router /auth/chat [post]
func AddChatRecord(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	p := new(models.ParamAddChatRecord)
	if err := ctx.ShouldBindJSON(p); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))
		return
	}
	code := service.AddChatRecord(p, user.(models.User).ID)
	if code != 0 {
		response.Response(ctx, http.StatusOK, code, nil, response.GetErrMsg(code))
		return
	}
	response.Success(ctx, gin.H{"tag": p}, "创建成功")
}

// DeleteChatRecord 删除聊天列表接口
// @Summary 删除聊天列表接口
// @Description 在聊天列表中删除一条记录，在路径中指定记录id，需要token
// @Tags 聊天列表相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path int true "记录id"
// @Success 200 {object} response.ResponseData
// @Router /auth/chat/{id} [delete]
func DeleteChatRecord(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var record models.ChatRecord
	mysql.RetrieveByID(&record, uint(id))
	//判断是否有权限
	user, _ := ctx.Get("user")
	if user.(models.User).ID != record.UserId {
		response.Response(ctx, http.StatusOK, response.Error, nil, "非法操作")
		return
	}
	err := mysql.Delete(&record)
	if err != nil {
		response.Response(ctx, http.StatusOK, response.CodeServerBusy, nil, "删除失败,该记录不存在")
		return
	}
	response.Success(ctx, gin.H{"record": record}, "删除成功")
}

// IsTop 置顶聊天接口
// @Summary 置顶聊天接口
// @Description 置顶，路径传id和is_top；is_top为0：不置顶，为1：置顶。需要token
// @Tags 聊天列表相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path int true "帖子id"
// @Param is_top path int true "状态"
// @Success 200 {object} response.ResponseData
// @Router /auth/chat/{id}/{is_top} [put]
func IsTop(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	isTop, _ := strconv.Atoi(ctx.Param("is_top"))
	var record models.ChatRecord
	mysql.RetrieveByID(&record, uint(id))

	//判断权限
	user, _ := ctx.Get("user")
	if user.(models.User).ID != record.UserId {
		response.Response(ctx, http.StatusOK, response.Error, nil, "非法操作")
		return
	}

	updateData := &models.ChatRecord{
		HasBg: isTop,
	}
	err := mysql.Update(record, updateData)
	if err != nil {
		response.Response(ctx, http.StatusOK, response.CodeServerBusy, nil, response.GetErrMsg(response.CodeServerBusy))
		return
	}
	response.Success(ctx, gin.H{"record": record}, "操作成功")
}
