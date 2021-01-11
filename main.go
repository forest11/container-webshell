package main

import (
	"github.com/astaxie/beego"
	_ "github.com/forest11/container-webshell/routers"
)

func main() {
	beego.Run()
}
