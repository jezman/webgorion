package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"

	_ "github.com/jezman/webgorion/routers"
)

func main() {
	beego.SetLevel(beego.LevelDebug)
	beego.SetLogger("file", `{"filename":"logs/webgorion.log"}`)
	sessionconf := &session.ManagerConfig{
		CookieName: "begoosessionID",
		Gclifetime: 3600,
	}
	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)
	go beego.GlobalSessions.GC()
	beego.Run()
}
