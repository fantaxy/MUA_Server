package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mua/controllers"
	"mua/models"
	"os"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true
	//强制每次删除重建表
	orm.RunSyncdb("default", false, true)
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/my", &controllers.MyEmotionController{})
	beego.Router("/my/delete", &controllers.MyEmotionController{}, "get:Delete")
	beego.Router("/my/modify", &controllers.MyEmotionController{}, "get:Modify")
	beego.Router("/data/all", &controllers.DataController{}, "get:All")

	// 附件处理
	os.Mkdir("emotions", os.ModePerm)
	beego.Router("/emotions/:all", &controllers.EmotionController{})
	// beego.SetStaticPath("/emotions", "emotions")

	beego.Run()
}
