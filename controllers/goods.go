package controllers

import (
	"encoding/json"
	"fourth-go-shFresh/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
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
	orm.Read(&goodsSku)

	//返回
	this.Data["goodsSku"] = goodsSku
	this.TplName = "detail.html"
}
