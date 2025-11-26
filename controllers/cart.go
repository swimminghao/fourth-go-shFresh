package controllers

import (
	"fourth-go-shFresh/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type CartController struct {
	beego.Controller
}

func (this *CartController) HandleAddCart() {
	//获取数据
	skuid, err1 := this.GetInt("skuid")
	count, err2 := this.GetInt("count")
	resp := make(map[string]interface{})
	if err1 != nil || err2 != nil {
		beego.Info("请求数据错误")
		resp["code"] = 1
		resp["msg"] = "传递数据不正确"
		return
	}
	beego.Info("skuid: ", skuid, "count：", count)
	//校验数据
	userName := this.GetSession("userName")
	if userName == nil {
		beego.Info("未登录状态")
		resp["code"] = 2
		resp["msg"] = "未登录"
		return
	}

	//处理数据
	conn, err := redis.Dial("tcp", "10.211.55.5:6379", redis.DialPassword("q123q123"))
	if err != nil {
		beego.Info("redis数据库连接err：", err)
		return
	}
	o := orm.NewOrm()
	user := models.User{Name: userName.(string)}
	o.Read(&user, "Name")
	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, count)
	//返回json数据

	resp["code"] = 5
	resp["msg"] = "ok"
	this.Data["json"] = resp
	this.ServeJSON()
}
