package controllers

import (
	"mua/models"

	"github.com/astaxie/beego"
)

type DataController struct {
	beego.Controller
}

func (this *DataController) All() {
	emotions, err := models.GetAllEmotions("", "", true)
	if err != nil {
		beego.Info(err)
		return
	}

	jsonStruct := make(map[string][]*models.Emotion)
	jsonStruct["data"] = emotions

	this.Data["json"] = jsonStruct
	this.ServeJson()
}
