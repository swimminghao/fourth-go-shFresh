package routers

import (
	"fourth-go-shFresh/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/user/*", beego.BeforeExec, filterFunc, false)
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")
	//激活用户
	beego.Router("/active", &controllers.UserController{}, "get:ActiveUser")
	//用户登录
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	//登录之后跳转请求
	beego.Router("/", &controllers.GoodsController{}, "get:ShowIndex")
	//退出登录
	beego.Router("/user/logout", &controllers.UserController{}, "get:GetLogout")
	//用户中心信息页
	beego.Router("/user/userCenterInfo", &controllers.UserController{}, "get:ShowUserCenterInfo")
	//用户中心订单页
	beego.Router("/user/userCenterOrder", &controllers.UserController{}, "get:ShowUserCenterOrder")
	//用户中心地址页
	beego.Router("/user/userCenterSite", &controllers.UserController{}, "get:ShowUserCenterSite;post:HandleUserCenterSite")
	//商品详情页面
	beego.Router("/goodsDetail", &controllers.GoodsController{}, "get:ShowGoodsDetail")
	//商品列表页
	beego.Router("/goodsList", &controllers.GoodsController{}, "get:ShowGoodsList")
	//商品搜索
	beego.Router("/goodsSearch", &controllers.GoodsController{}, "post:HandleSearch")
	beego.Router("/user/addCart", &controllers.CartController{}, "post:HandleAddCart")

}

var filterFunc = func(ctx *context.Context) {

	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
