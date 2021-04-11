package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hackathon/service/websocket"
)

func WebSocket(ctx *gin.Context) {

	WS := new(websocket.Ws)

	fmt.Println("haha")

	cli,_ := WS.OnOpen(ctx)

	fmt.Println("1:",cli ,"2:", cli.WsClient)

	cli.OnMessage(ctx)

}
