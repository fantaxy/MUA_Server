package routers

import (
	"github.com/astaxie/beego"
	"mua/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
