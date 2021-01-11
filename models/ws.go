package models

import (
	"net"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"

	"encoding/json"
	"log"
)

// websocket 连接器
type Connection struct {
	Ws *websocket.Conn

	// 发送信息的缓冲 channel
	Send chan []byte
}

type toJson struct {
	Data string `json:"data"`
}

// 处理ws请求
func (c *Connection) ConnHandle(conn net.Conn) {
	for {
		_, msg, err := c.Ws.ReadMessage()
		if err != nil {
			// 判断是不是超时
			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() {
					log.Printf("ReadMessage timeout remote: %v\n", conn.RemoteAddr())
					return
				}
			}
			log.Printf("ReadMessage remote:%v error: %v \n", conn.RemoteAddr(), err)
			return
		}
		conn.Write(msg)
	}
}

// 发送数据
func (c *Connection) SendResp(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		conn.Read(buf)
		jData := toJson{Data: string(buf)}
		data, _ := json.Marshal(jData)
		err := c.Ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("send msg faild:  %v\n", err)
			return
		}
		beego.Info(string(buf))
	}
}
