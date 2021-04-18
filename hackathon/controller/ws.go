package controller

import (
	"github.com/gin-gonic/gin"
	"hackathon/models"
	"hackathon/response"
	"hackathon/service"
	"hackathon/service/websocket"
	"net/http"
)

func WebSocket(ctx *gin.Context) {

	WS := new(websocket.Ws)
	cli,_ := WS.OnOpen(ctx)
	cli.OnMessage(ctx)

}

type BindDate struct {

	From      int       `json:"from" form:"from"`
	To        int       `json:"to" form:"to"`
	Last    int `json:"last"`

}

type ChatData struct {

	From []models.ChatHistory
	To   []models.ChatHistory

}

func ChatHistory(ctx *gin.Context){

	var history = models.ChatHistory{}
	var bindDate = BindDate{}

	if err := ctx.ShouldBind(&bindDate); err != nil {
		response.Response(ctx, http.StatusOK, response.CodeParamError, nil, response.GetErrMsg(response.CodeParamError))

	}


	history.To =bindDate.To
	history.From = bindDate.From

	from , to := service.History(history, bindDate.Last)

	data := ChatData{From: from , To: to}

	ctx.JSON(http.StatusOK,gin.H{

		"code" : 200 ,
		"data" : data,

	})


}





