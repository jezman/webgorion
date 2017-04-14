package main

import (
	"github.com/astaxie/beego"

	_ "github.com/jezman/webgorion/routers"
)

func main() {
	beego.SetLevel(beego.LevelDebug)
	beego.SetLogger("file", `{"filename":"logs/webgorion.log"}`)
	beego.Run()
}
