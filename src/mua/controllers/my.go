package controllers

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"mua/models"
	// "os"
	"path"
)

type MyEmotionController struct {
	beego.Controller
}

func (this *MyEmotionController) Get() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.Data["IsMy"] = true
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	uname := currentUser(this.Ctx)
	emotions, err := models.GetAllEmotions(uname, "", true)
	if err != nil {
		beego.Error(err)
	}
	this.Data["Emotions"] = emotions
	this.TplNames = "my.html"
}

func (this *MyEmotionController) Post() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	tags := this.Input().Get("tags")
	eid := this.Input().Get("eid")
	uname := currentUser(this.Ctx)

	if len(eid) == 0 {
		//上传表情
		f, fh, err := this.GetFile("image")
		if err != nil {
			beego.Error(err)
			this.Redirect("/my", 302)
			return
		}

		//读取前500字节生成MD5
		b := make([]byte, 500)
		f.Read(b)
		defer f.Close()

		var image string
		var md5Name string
		if fh != nil {
			image = fh.Filename
			beego.Info(image)
			md5Name = fmt.Sprintf("%x", md5.Sum(b))
			beego.Info(md5Name)
		}
		//查重
		if models.CheckDuplicate(md5Name) {
			this.Redirect("/my", 302)
			beego.Warning("Duplicate image: " + image)
			return
		}
		//保存到文件
		err = this.SaveToFile("image", path.Join("emotions", md5Name))
		if err != nil {
			beego.Error(err)
		}

		//保存到数据库
		err = models.AddEmotion(md5Name, uname, tags)
	} else {
		//修改表情
		err := models.ModifyEmotion(eid, uname, tags)
		if err != nil {
			beego.Error(err)
			this.Redirect("/my", 302)
			return
		}
	}

	this.Redirect("/my", 302)
}

func (this *MyEmotionController) Delete() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	err := models.DeleteEmotion(this.Input().Get("eid"))
	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/my", 302)
}

func (this *MyEmotionController) Modify() {
	this.TplNames = "item_modify.html"

	eid := this.Input().Get("eid")
	emotion, err := models.GetEmotion(eid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	this.Data["Emotion"] = emotion
	this.Data["Eid"] = eid
	this.Data["IsLogin"] = true
}
