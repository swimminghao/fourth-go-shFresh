package main

import (
	_ "fourth-go-shFresh/routers"
	_ "fourth-go-shFresh/models"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

