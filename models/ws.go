package models

import (
	"encoding/json"
	"io"
	"log"
	"net"

	"github.com/gorilla/websocket"
)

// websocket 连接器
type Connection struct {
	Ws *websocket.Conn

	// 发送信息的缓冲 channel
	Send chan []byte
}

// json数据
type toJson struct {
	Data string `json:"data"`
}

// 处理ws请求
func (c *Connection) ConnHandle(conn net.Conn) {
	for {
		_, msg, err := c.Ws.ReadMessage()
		if err != nil {
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
		_, err := conn.Read(buf)
		if err == io.EOF { // 与docker断开连接, 需要关闭websocket连接
			c.Ws.Close()
			return
		}
		jData := toJson{Data: string(buf)}
		data, _ := json.Marshal(jData)
		err = c.Ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("send msg faild:  %v\n", err)
			return
		}
	}
}
