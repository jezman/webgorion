package controllers

import "github.com/astaxie/beego"

type SessionController struct {
	beego.Controller
}

func (c *SessionController) Login() {
	c.Data["Title"] = "Login page"
	c.TplName = "login.html"
	c.Layout = "login.html"
}
