package controllers

import (
	"encoding/json"
	"fourth-go-shFresh/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type GoodsController struct {
	beego.Controller
}

//
func GetUser(this *beego.Controller) string {
	userName := this.GetSession("userName")
	if userName != nil {
		this.Data["userName"] = userName.(string)
		return userName.(string)
	} else {
		this.Data["userName"] = ""
	}
	return ""
}
func (this *GoodsController) ShowIndex() {
	GetUser(&this.Controller)
	orm := orm.NewOrm()
	//获取类型数据
	var goodsTypes []models.GoodsType
	orm.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes
	beego.Info("goodsTypes:", goodsTypes)
	//获取轮播图数据
	var indexGoods []models.IndexGoodsBanner
	orm.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&indexGoods)
	this.Data["indexGoodsBanner"] = indexGoods
	beego.Info("indexGoodsBanner:", indexGoods)

	//获取促销商品数据
	var promotionGoods []models.IndexPromotionBanner
	orm.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionGoods)
	this.Data["promotionGoods"] = promotionGoods
	//获取展示商品数据
	goods := make([]map[string]interface{}, len(goodsTypes))
	//向切片interface类型中添加类型数据
	for index, value := range goodsTypes {
		//获取对应类型的首页展示商品
		temp := make(map[string]interface{})
		temp["type"] = value
		goods[index] = temp
	}
	//商品数据

	for _, value := range goods {
		var textGoods []models.IndexTypeGoodsBanner
		var imgGoods []models.IndexTypeGoodsBanner
		//获取文字商品数据
		orm.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").OrderBy("Index").Filter("DisplayType", 0).Filter("GoodsType", value["type"]).All(&textGoods)
		//获取图片商品数据
		orm.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").OrderBy("Index").Filter("DisplayType", 1).Filter("GoodsType", value["type"]).All(&imgGoods)

		value["textGoods"] = textGoods
		value["imgGoods"] = imgGoods
	}
	this.Data["goods"] = goods
	jsonString := ToJSONStringSafe(goods)
	logs.Info("goods:", len(goods), jsonString)
	this.TplName = "index.html"

}

// 安全转换 - 出错时返回空对象
func ToJSONStringSafe(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
func ShowLaout(this *beego.Controller) {
	//查询类型
	orm := orm.NewOrm()
	var types []models.GoodsType
	orm.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types
	//获取用户信息
	GetUser(this)
	//指定layout
	this.Layout = "goodsLayout.html"
}

//展示商品详情
func (this *GoodsController) ShowGoodsDetail() {
	//获取
	id, err := this.GetInt("id")

	//校验
	if err != nil {
		beego.Error("浏览器请求错误")
		this.Redirect("/", 302)
		return
	}
	//处理
	orm := orm.NewOrm()
	var goodsSku models.GoodsSKU
	goodsSku.Id = id
	//orm.Read(&goodsSku)
	orm.QueryTable("GoodsSKU").RelatedSel("GoodsType", "Goods").Filter("Id", id).One(&goodsSku)
	//获取同类型时间靠前的两条商品数据
	var goodsNew []models.GoodsSKU
	orm.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType", goodsSku.GoodsType).OrderBy("Time").Limit(2, 0).All(&goodsNew)
	//判断用户是否登录
	userName := this.GetSession("userName")
	if userName != nil {
		user := models.User{Name: userName.(string)}
		orm.Read(&user, "Name")
		//添加历史浏览记录,用redis存储
		conn, err := redis.Dial("tcp", "10.211.55.5:6379", redis.DialPassword("q123q123"))
		defer conn.Close()
		if err != nil {
			beego.Info("redis连接错误")
		}
		//把以前相同商品的历史浏览记录删除
		conn.Do("lrem", "history_"+strconv.Itoa(user.Id), 0, id)
		//添加新的商品浏览记录
		conn.Do("lpush", "history_"+strconv.Itoa(user.Id), id)
	}
	//添加历史记录
	//返回
	this.Data["goodsSku"] = goodsSku
	this.Data["goodsNew"] = goodsNew
	ShowLaout(&this.Controller)
	this.TplName = "detail.html"
}

//展示商品列表页
func (this *GoodsController) ShowGoodsList() {
	id, err := this.GetInt("typeId")
	if err != nil {
		beego.Info("请求路径错误")
		this.Redirect("/", 302)
		return
	}
	//处理数据
	ShowLaout(&this.Controller)
	//获取新品`
	o := orm.NewOrm()
	var goodsNew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Time").Limit(2, 0).All(&goodsNew)
	this.Data["goodsNew"] = goodsNew

	// 分页参数
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil || pageIndex < 1 {
		pageIndex = 1
	}
	pageSize := 3

	// 获取总记录数
	count, _ := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).Count()

	// 使用 CreatePagination 创建分页对象
	pagination := CreatePagination(count, pageIndex, pageSize)

	// 获取当前页数据
	var goods []models.GoodsSKU
	query := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id)

	// 排序处理
	sort := this.GetString("sort")
	switch sort {
	case "price":
		query = query.OrderBy("Price")
		this.Data["sort"] = "price"
	case "sale":
		query = query.OrderBy("Sales")
		this.Data["sort"] = "sale"
	default:
		this.Data["sort"] = ""
	}

	// 计算偏移量并查询数据
	offset := (pagination.PageIndex - 1) * pagination.PageSize
	query.Limit(pagination.PageSize, offset).All(&goods)

	// 设置模板数据
	this.Data["goods"] = goods
	this.Data["pagination"] = pagination // 传递整个分页对象
	this.Data["typeId"] = id

	this.TplName = "list.html"
}

//商品搜索
func (this *GoodsController) HandleSearch() {
	goodsName := this.GetString("goodsName")
	o := orm.NewOrm()
	var goods []models.GoodsSKU
	if goodsName == "" {
		o.QueryTable("GoodsSKU").All(&goods)
		this.Data["goods"] = goods
		ShowLaout(&this.Controller)
		this.TplName = "search.html"
		return
	}
	o.QueryTable("GoodsSKU").Filter("Name__icontains", goodsName).All(&goods)
	this.Data["goods"] = goods
	beego.Info("goodsSearch: ", goods)
	ShowLaout(&this.Controller)
	this.TplName = "search.html"
}
