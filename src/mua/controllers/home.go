package controllers

import (
	"github.com/astaxie/beego"
	"mua/models"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["IsHome"] = true
	this.TplNames = "home.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	dict, err := models.GetAllSeries()
	if err != nil {
		beego.Info(err)
		return
	}
	this.Data["emotionDict"] = dict
}
