package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hackathon/core"
	"hackathon/dao/mysql"
	"hackathon/models"
	"log"
	"strconv"
	"strings"
	"time"
)

type Ws struct {
	WsClient *core.Client
}



// onOpen 基本不需要做什么
func (w *Ws) OnOpen(context *gin.Context) (*Ws, bool) {
	if client, ok := (&core.Client{}).OnOpen(context); ok {
		w.WsClient = client
		go w.WsClient.Heartbeat(w.OnClose) // 一旦握手+协议升级成功，就为每一个连接开启一个自动化的隐式心跳检测包
		return w, true
	} else {
		return nil, false
	}
}

// OnMessage 处理业务消息
func (w *Ws) OnMessage(context *gin.Context) {
	go w.WsClient.ReadPump(func(messageType int, receivedData []byte) {
		//参数说明
		//messageType 消息类型，1=文本
		//receivedData 服务器接收到客户端（例如js客户端）发来的的数据，[]byte 格式

		//tempMsg := "服务器已经收到了你的消息==>" + string(receivedData)
		// 回复客户端已经收到消息;
		//if err := w.WsClient.Conn.WriteMessage(messageType, []byte(tempMsg)); err != nil {
		//	log.Fatal(err.Error())
		//}
		TempStr := string(receivedData)
		Data := strings.Split(TempStr,":")
		To ,err:= strconv.Atoi(Data[0])
		var Msg string
		if err == nil {
			Msg = Data[1]
		}else{
			To = w.WsClient.Uid
			Msg = err.Error()
		}

		//for key, message := range Hub{
		//
		//	if message.to == w.WsClient.Uid {
		//
		//		w.WsClient.Conn.WriteMessage(messageType,[]byte(message.message))
		//
		//		Hub = append(Hub[:key], Hub[key+1:]...)
		//
		//	}
		//
		//
		//}
		if _,ok := w.WsClient.Hub.ClientsId[To];ok {

			message :=models.ChatHistory{From: w.WsClient.Uid, To: To,Msg: Msg,CreatedAt: time.Now()}
			mysql.Create(&message)
			err := w.WsClient.Hub.ClientsId[To].Conn.WriteMessage(messageType, []byte(strconv.Itoa(w.WsClient.Uid)+":"+Msg))
			if err != nil {

				fmt.Println(err.Error())

			}
		}else {

			message :=models.ChatHistory{From: w.WsClient.Uid, To: To,Msg: Msg,CreatedAt: time.Now()}
			mysql.Create(&message)

			w.WsClient.Producer.Send(strconv.Itoa(To) , strconv.Itoa(w.WsClient.Uid)+":"+Msg)
			w.WsClient.Producer.Close()
		}

	}, w.OnError, w.OnClose)
}

// OnError 客户端与服务端在消息交互过程中发生错误回调函数
func (w *Ws) OnError(err error) {
	panic(err)

	//fmt.Printf("远端掉线、卡死、刷新浏览器等会触发该错误: %v\n", err.Error())
}

// OnClose 客户端关闭回调，发生onError回调以后会继续回调该函数
func (w *Ws) OnClose() {
	w.WsClient.Hub.UnRegister <- w.WsClient // 向hub管道投递一条注销消息，有hub中心负责关闭连接、删除在线数据
}

//获取在线的全部客户端
func (w *Ws) GetOnlineClients() {

	fmt.Printf("在线客户端数量：%d\n", len(w.WsClient.Hub.Clients))
}

// 向全部在线客户端广播消息
func (w *Ws) BroadcastMsg(sendMsg string) {

	for onlineClient := range w.WsClient.Hub.Clients {

		// 每次向客户端写入消息命令（WriteMessage）之前必须设置超时时间
		if err := onlineClient.Conn.SetWriteDeadline(time.Now().Add(w.WsClient.WriteDeadline * time.Second)); err != nil {
			log.Fatal(err.Error())
		}
		//获取每一个在线的客户端，向远端发送消息
		if err := onlineClient.Conn.WriteMessage(websocket.TextMessage, []byte(sendMsg)); err != nil {
			log.Fatal(err.Error())
		}
	}
}

