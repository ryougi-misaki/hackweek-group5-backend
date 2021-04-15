package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hackathon/service/websocket"
)

func WebSocket(ctx *gin.Context) {

	fmt.Println("haha1")
	WS := new(websocket.Ws)

	fmt.Println("haha")

	cli,_ := WS.OnOpen(ctx)

	cli.OnMessage(ctx)

	//cli.BroadcastMsg("新成员加入")

	cli.GetOnlineClients()

}
