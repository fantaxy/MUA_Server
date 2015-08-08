package controllers

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"mua/models"
	"net/http"
	"os"
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
	w := this.Ctx.ResponseWriter

	if len(eid) == 0 {
		//上传表情
		files, err := this.GetFiles("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
			this.Redirect("/my", 302)
			return
		}
		for i, _ := range files {
			fh := files[i]
			f, err := fh.Open()
			defer f.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//读取前500字节生成MD5
			b := make([]byte, 500)
			f.Read(b)
			f.Seek(0, 0)

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
			//create destination file making sure the path is writeable.
			dst, err := os.Create("emotions/" + md5Name)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, f); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//保存到数据库
			err = models.AddEmotion(md5Name, uname, tags)
		}
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
