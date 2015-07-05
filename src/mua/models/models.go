package models

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	//加下划线表示只执行初始化函数
	_ "github.com/mattn/go-sqlite3"
)

type Emotion struct {
	Id         int64
	Name       string
	Uploader   string
	UploadTime time.Time `orm:"index"`
	Downloads  int64
	Series     string
	Tags       string
}

type Series struct {
	Id          int64
	Name        string
	CreateTime  time.Time `orm:"index"`
	Description string
	ThumbName   string
	Tag         string
}

type Tag struct {
	Id         int64
	Name       string `orm:"index;unique"`
	CreatedBy  string
	CreateTime time.Time `orm:"index"`
}

const (
	_DB_NAME        = "data/mua.db"
	_SQLITE3_DRIVER = "sqlite3"
)

func RegisterDB() {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	orm.RegisterModel(new(Emotion), new(Series), new(Tag))
	orm.RegisterDriver((_SQLITE3_DRIVER), orm.DR_Sqlite)
	//必须要有注册为default的数据库
	//最大连接数10
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

func AddEmotion(name, uploader, tags string) error {
	// 处理标签
	AddTags(uploader, tags)
	tags = "^" + strings.Join(strings.Split(tags, " "), "$^") + "$"

	o := orm.NewOrm()

	emotion := &Emotion{
		Name:       name,
		Uploader:   uploader,
		UploadTime: time.Now(),
		Tags:       tags,
	}
	_, err := o.Insert(emotion)
	if err != nil {
		return err
	}

	return err
}

func AddTags(uploader string, tags string) error {
	o := orm.NewOrm()
	var err error
	for _, v := range strings.Split(tags, " ") {
		tag := &Tag{
			Name:      v,
			CreatedBy: uploader,
			CreateTime:time.Now()
		}
		_, err = o.Insert(tag)
	}
	return err
}

func GetAllSeries() (dict map[string][]*Emotion, err error) {
	dict = make(map[string][]*Emotion)
	tags := make([]*Tag, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("tag")
	qs.OrderBy("-profile__age", "profile")
	qs.All(&tags)

	for i := range tags {
		tagName := tags[i].Name
		emotions, err := GetAllEmotions("", tagName, false)
		// beego.Debug("tag: ", tagName, " emotions: ", emotions)
		if err != nil || len(emotions) == 0 {
			continue
		}
		if dict[tagName] == nil {
			dict[tagName] = emotions
		} else {
			dict[tagName] = append(dict[tagName], emotions...)
		}
	}
	return dict, err
}

func GetAllEmotions(uploader string, tag string, isDesc bool) (emotions []*Emotion, err error) {
	o := orm.NewOrm()

	emotions = make([]*Emotion, 0)

	qs := o.QueryTable("emotion")
	if len(uploader) > 0 {
		qs = qs.Filter("uploader", uploader)
	}
	if len(tag) > 0 {
		qs = qs.Filter("tags__contains", "^"+tag+"$")
	}
	if isDesc {
		_, err = qs.OrderBy("-upload_time").All(&emotions)

	} else {
		_, err = qs.All(&emotions)
	}
	for i := range emotions {
		emotions[i].Tags = strings.Replace(strings.Replace(
			emotions[i].Tags, "$", " ", -1), "^", "", -1)
	}

	return emotions, err
}

func GetEmotion(id string) (Emotion, error) {
	o := orm.NewOrm()

	eid, err := strconv.ParseInt(id, 10, 64)
	emotion := Emotion{Id: eid}

	err = o.Read(&emotion)
	emotion.Tags = strings.Replace(strings.Replace(
		emotion.Tags, "$", " ", -1), "^", "", -1)

	return emotion, err
}

func CheckDuplicate(name string) bool {
	o := orm.NewOrm()

	qs := o.QueryTable("emotion")
	if len(name) > 0 {
		qs = qs.Filter("name", name)
	}
	count, _ := qs.Count()
	if count > 0 {
		return true
	} else {
		return false
	}
}

func DeleteEmotion(id string) error {
	eid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	emotion := &Emotion{Id: eid}
	_, err = o.Delete(emotion)
	return err
}

func ModifyEmotion(id, uname, tags string) error {
	eid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	escapedTags := "^" + strings.Join(strings.Split(tags, " "), "$^") + "$"
	emotion := &Emotion{Id: eid, Tags: escapedTags}
	_, err = o.Update(emotion, "Tags")
	if err == nil {
		err = AddTags(uname, tags)
	}

	return err
}
