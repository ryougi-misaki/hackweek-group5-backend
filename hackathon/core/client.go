package core

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hackathon/rabbitMq"
	"hackathon/variable"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	Hub                *Hub            // 负责处理客户端注册、注销、在线管理
	Conn               *websocket.Conn // 一个ws连接
	Send               chan []byte     // 一个ws连接存储自己的消息管道
	PingPeriod         time.Duration
	ReadDeadline       time.Duration
	WriteDeadline      time.Duration
	HeartbeatFailTimes int

	//rabbitMq

	Consume *rabbitMq.Consume

	Producer *rabbitMq.Producer

	// 自己的地址
	ID interface{}
	//用户的id
	Uid int
	//目标id
	To int
}

// 处理握手+协议升级
func (c *Client) OnOpen(context *gin.Context) (*Client, bool) {
	// 1.升级连接,从http--->websocket
	defer func() {
		err := recover()
		if err != nil {

			if _, ok := err.(error); ok {
			}

		}
	}()
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  2000,
		WriteBufferSize: 2000,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 2.将http协议升级到websocket协议.初始化一个有效的websocket长连接客户端
	if wsConn, err := upGrader.Upgrade(context.Writer, context.Request, nil); err != nil {
		return nil, false
	} else {
		if wsHub, ok := variable.WebsocketHub.(*Hub); ok {
			c.Hub = wsHub
		}
		c.Conn = wsConn
		c.Send = make(chan []byte, 20480)
		c.PingPeriod = time.Second * 30
		c.ReadDeadline = time.Second * 0
		c.WriteDeadline = time.Second * 35
		c.ID = c
		c.Uid, _ = strconv.Atoi(context.Param("id"))
		c.Consume, _ = rabbitMq.CreateConsumer(strconv.Itoa(c.Uid))
		c.Producer, _ = rabbitMq.CreateProducer(strconv.Itoa(c.Uid))
		c.Conn.SetReadLimit(65535) // 设置最大读取长度
		c.Hub.Register <- c

		//注册后,收取积攒的消息

		go c.Consume.Received(func(receivedData string) {

			c.Conn.WriteMessage(websocket.TextMessage, []byte(receivedData))

		})

		return c, true
	}

}

// 主要功能主要是实时接收消息
func (c *Client) ReadPump(callbackOnMessage func(messageType int, receivedData []byte), callbackOnError func(err error), callbackOnClose func()) {
	// 回调 onclose 事件
	defer func() {
		err := recover()
		if err != nil {
			if realErr, isOk := err.(error); isOk {
				panic(realErr)
			}
		}
		callbackOnClose()
	}()

	// OnMessage事件
	for {
		mt, bReceivedData, err := c.Conn.ReadMessage()
		if err == nil {
			if err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteDeadline)); err != nil {
				panic(err)
			}
			callbackOnMessage(mt, bReceivedData)
		} else {
			// OnError事件
			//callbackOnError(err)
			break
		}
	}
}

// 按照websocket标准协议实现隐式心跳,Server端向Client远端发送ping格式数据包,浏览器收到ping标准格式，自动将消息原路返回给服务器
func (c *Client) Heartbeat(callbackClose func()) {
	//  1. 设置一个时钟，周期性的向client远端发送心跳数据包
	ticker := time.NewTicker(c.PingPeriod)
	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				panic(val)
			}
		}
		ticker.Stop()   // 停止该client的心跳检测
		callbackClose() // 注销 client
	}()
	//2.浏览器收到服务器的ping格式消息，会自动响应pong消息，将服务器消息原路返回过来
	if c.ReadDeadline == 0 {
		_ = c.Conn.SetReadDeadline(time.Time{})
	}
	c.Conn.SetPongHandler(func(receivedPong string) error {
		if c.ReadDeadline > time.Nanosecond {
			_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
		} else {
			_ = c.Conn.SetReadDeadline(time.Time{})
		}
		//fmt.Println("浏览器收到ping标准格式，自动将消息原路返回给服务器：", received_pong)  // 接受到的消息叫做pong，实际上就是服务器发送出去的ping数据包
		return nil
	})
	//3.自动心跳数据
	for {
		select {
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(c.WriteDeadline))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte(variable.WebsocketServerPingMsg)); err != nil {
				c.HeartbeatFailTimes++
				if c.HeartbeatFailTimes > 4 {
					callbackClose()
					return
				}
			} else {
				if c.HeartbeatFailTimes > 0 {
					c.HeartbeatFailTimes--
				}
			}
		}
	}
}
