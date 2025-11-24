package controllers

import (
	"encoding/base64"
	"fourth-go-shFresh/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"regexp"
	"strconv"
)

type UserController struct {
	beego.Controller
}

//显示注册页面
func (this *UserController) ShowReg() {
	this.TplName = "register.html"
}

//处理注册数据
func (this *UserController) HandleReg() {
	//1.获取数据
	userName := this.GetString("user_name")
	pwd := this.GetString("pwd")
	cpwd := this.GetString("cpwd")
	email := this.GetString("email")
	//2.校验数据
	if userName == "" || pwd == "" || cpwd == "" || email == "" {
		this.Data["errmsg"] = "数据不完整，请重新检查数据~"
		this.TplName = "register.html"
		return
	}
	if pwd != cpwd {
		this.Data["errmsg"] = "两次密码不一致，请重新输入！"
		this.TplName = "register.html"
		return
	}
	regexp, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := regexp.FindString(email)
	if res == "" {
		this.Data["errmsg"] = "邮箱格式不正确，请重新输入！"
		this.TplName = "register.html"
		return
	}
	//3.处理数据
	orm := orm.NewOrm()
	user := models.User{Name: userName, PassWord: pwd, Email: email}
	_, err := orm.Insert(&user)
	if err != nil {
		this.Data["errmsg"] = "注册失败，请更新数据重新注册！"
		this.TplName = "register.html"
		return
	}
	//发送邮箱
	emailConfig := `{"username":"18709211491@163.com","password":"VMTkVKmAdyEpqzr7","host":"smtp.163.com","port":25}`
	emailConn := utils.NewEMail(emailConfig)
	emailConn.From = "18709211491@163.com"
	emailConn.To = []string{email}
	emailConn.Subject = "天天生鲜用户注册"
	//注意这里是发送给用户的激活请求地址
	activeUrl := "http://10.211.55.5:8080/active?id=" + strconv.Itoa(user.Id)
	s := `<a href="` + activeUrl + `">激活链接</a>`
	beego.Info("激活链接：", s)
	emailConn.HTML = s
	err = emailConn.Send()
	if err != nil {
		beego.Info("激活邮件发送失败", err)
	}

	//4.返回数据
	this.Ctx.WriteString("注册成功，请激活用户。")
}

//激活处理
func (this *UserController) ActiveUser() {
	//获取数据
	id, err := this.GetInt("id")
	//校验数据
	if err != nil {
		this.Data["errmsg"] = "要激活的用户不存在"
		this.TplName = "register.html"
		return
	}
	//处理数据
	//更新操作
	orm := orm.NewOrm()
	user := models.User{Id: id}
	err = orm.Read(&user)
	if err != nil {
		this.Data["errmsg"] = "要激活的用户不存在"
		this.TplName = "register.html"
		return
	}
	user.Active = true
	orm.Update(&user)
	//返回数据
	this.Redirect("/login", 302)
}

//展示登录页面
func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	//解码
	temp, _ := base64.StdEncoding.DecodeString(userName)
	if string(temp) == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}

	this.TplName = "login.html"
}

//处理登录业务
func (this *UserController) HandleLogin() {
	userName := this.GetString("username")
	pwd := this.GetString("pwd")
	beego.Info(userName, pwd)
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "登录数据不完整，请重新登录"
		this.TplName = "login.html"
		return
	}
	orm := orm.NewOrm()
	user := models.User{Name: userName}
	err := orm.Read(&user, "Name")
	if err != nil {
		beego.Info("错误：", err)
		this.Data["errmsg"] = "用户名或密码不存在，请重新登录"
		this.TplName = "login.html"
		return
	}
	if pwd != user.PassWord {
		this.Data["errmsg"] = "用户名或密码不正确，请重新登录"
		this.TplName = "login.html"
		return
	}
	if user.Active != true {
		this.Data["errmsg"] = "用户名未激活，请前往邮箱激活"
		this.TplName = "login.html"
		return
	}

	remember := this.GetString("remember")
	if remember == "on" {
		temp := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName", temp, 24*3600*30)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)

	}
	//跳转到首页,
	/*
		1.首页的简单显示实现
		2.登录判断（路由过滤器）
		3.首页显示
		4.三个页面
			视图布局
			添加地址页（如何让页面只显示一个地址）
			用户中心信息页显示
	*/
	this.SetSession("userName", userName)
	//this.Ctx.WriteString("登录成功！")
	this.Redirect("/", 302)
}

//退出登录
func (this *UserController) GetLogout() {
	this.DelSession("userName")
	this.Redirect("/login", 302)
}

//展示用户中心信息页面
func (this *UserController) ShowUserCenterInfo() {
	userName := GetUser(&this.Controller)
	this.Data["userName"] = userName
	orm := orm.NewOrm()
	var addr models.Address
	orm.QueryTable("Address").RelatedSel("User").Filter("User__Name", userName).Filter("IsDefault", true).One(&addr)
	if addr.Id == 0 {
		this.Data["addr"] = ""
	} else {
		this.Data["addr"] = addr
	}
	beego.Info("用户详情bo:", addr)
	this.Layout = "userCenterLayout.html"
	this.TplName = "user_center_info.html"
}

//展示用户订单页
func (this *UserController) ShowUserCenterOrder() {
	GetUser(&this.Controller)

	this.Layout = "userCenterLayout.html"
	this.TplName = "user_center_order.html"
}

//展示用户地址页
func (this *UserController) ShowUserCenterSite() {
	userName := GetUser(&this.Controller)
	//this.Data["userName"] = userName
	orm := orm.NewOrm()
	var addr models.Address
	orm.QueryTable("Address").RelatedSel("User").Filter("User__Name", userName).Filter("IsDefault", true).One(&addr)
	this.Data["addr"] = addr
	this.Layout = "userCenterLayout.html"
	this.TplName = "user_center_site.html"
}

func (this *UserController) HandleUserCenterSite() {
	receiver := this.GetString("receiver")
	zipcode := this.GetString("zipcode")
	phone := this.GetString("phone")
	addr := this.GetString("addr")
	if addr == "" || receiver == "" || phone == "" || zipcode == "" {
		beego.Info("数据不完整")
		this.Redirect("/user/userCenterSite", 302)
	}
	orm := orm.NewOrm()
	addrUser := models.Address{IsDefault: true}
	err := orm.Read(&addrUser, "IsDefault")
	if err == nil {
		addrUser.IsDefault = false
		orm.Update(&addrUser)
	}
	userName := this.GetSession("userName")
	user := models.User{Name: userName.(string)}
	err = orm.Read(&user, "Name")
	addrUserDefault := models.Address{Receiver: receiver, Addr: addr, ZipCode: zipcode, Phone: phone, User: &user}
	orm.Insert(&addrUserDefault)

	this.Redirect("/user/userCenterSite", 302)
}
