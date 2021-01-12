package controllers

import (
	"fmt"
	"log"
	"net"

	"github.com/astaxie/beego"
	"github.com/forest11/container-webshell/handler"
	"github.com/forest11/container-webshell/models"
	"github.com/gorilla/websocket"
)

type Wscontroller struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{}

func (self *Wscontroller) Get() {
	host := self.Input().Get("h")
	port := self.Input().Get("p")
	containerid := self.Input().Get(("containers_id"))
	rows := self.Input().Get("rows")
	cols := self.Input().Get("cols")

	execid, err := handler.GetDockerExecId(host, port, containerid)
	if err != nil {
		beego.Error(err)
	}

	log.Println("execid is", execid)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		beego.Error(err)
	}
	defer conn.Close()

	data := "{\"Tty\":true}"
	_, err = conn.Write([]byte(fmt.Sprintf("POST /exec/%s/start HTTP/1.1\r\nHost: %s\r\nContent-Type: application/json\r\nContent-Length: %s\r\n\r\n%s",
		execid, fmt.Sprintf("%s:%s", host, port), fmt.Sprint(len([]byte(data))), data)))

	if err != nil {
		beego.Error(err)
	}

	handler.ResizeContainer(host, port, execid, rows, cols)
	ws, err := upgrader.Upgrade(self.Ctx.ResponseWriter, self.Ctx.Request, nil)
	c := &models.Connection{Ws: ws, Send: make(chan []byte, 512)}
	go c.ConnHandle(conn)
	c.SendResp(conn)
}
