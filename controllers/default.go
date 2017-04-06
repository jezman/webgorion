package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/denisenkom/go-mssqldb"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Redirect("/acs/summary", 302)
}
