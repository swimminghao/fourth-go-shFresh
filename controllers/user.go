package controllers

import (
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowReg() {
	this.TplName = "register.html"
}
