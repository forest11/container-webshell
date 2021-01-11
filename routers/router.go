package routers

import (
	"github.com/astaxie/beego"
	"github.com/forest11/container-webshell/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.Wscontroller{})
}
