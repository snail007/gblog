package admin

import (
	"github.com/snail007/gmc/util/cast"
	gmap "github.com/snail007/gmc/util/map"
	"gblog/model"
	"time"
)

type User struct {
	Admin
}

func (this *User) Profile() {
	if this.Ctx.IsAJAX() {
		nickname := this.Ctx.POST("nickname")
		this.User["nickname"] = nickname
		data := gmap.M{
			"nickname":    nickname,
			"update_time": time.Now().Unix(),
		}
		_, err := model.User.UpdateByIDs([]string{this.User["user_id"]}, data)
		this.StopE(err, func() {
			this._JSONFail(err.Error())
		}, func() {
			this.Session.Set("admin", this.User)
			this._JSONSuccess("", nil, this.url("/user/profile"))
		})
	} else {
		this.View.Layout("admin/page").Render("admin/user/profile")
	}
}

func (this *User) Password() {
	if this.Ctx.IsAJAX() {
		password := this.Ctx.POST("password")
		password1 := this.Ctx.POST("password1")
		password2 := this.Ctx.POST("password2")
		if password == "" || password1 == "" || (password1 != password2) {
			this._JSONFail("信息不完整")
		}
		password0 := ""
		dbUser, err := model.User.GetByID(this.User["user_id"])
		if err != nil || len(dbUser) == 0 {
			this._JSONFail("用户不存在")
		}
		if dbUser["password"] != model.EncodePassword(password) {
			this._JSONFail("当前密码错误")
		}
		password0 = model.EncodePassword(password1)
		data := gmap.M{
			"password":    password0,
			"update_time": time.Now().Unix(),
		}
		_, err = model.User.UpdateByIDs([]string{this.User["user_id"]}, data)
		this.StopE(err, func() {
			this._JSONFail(err.Error())
		}, func() {
			this._JSONSuccess("", nil, this.url("/user/password"))
		})
	} else {
		this.View.Layout("admin/page").Render("admin/user/password")
	}
}

func (this *User) Add() {
	this.View.Layout("admin/form").Render("admin/user/form")
}

func (this *User) Create() {
	username := this.Ctx.POST("username")
	nickname := this.Ctx.POST("nickname")
	password := this.Ctx.POST("password")
	password1 := this.Ctx.POST("password1")
	if username == "" || password == "" || password != password1 {
		this._JSONFail("信息不完整")
	}
	dbUser, err := model.User.GetBy(gmap.M{"username": username})
	if err != nil || len(dbUser) > 0 {
		this._JSONFail("用户已经存在")
	}
	now := time.Now().Unix()
	data := gmap.M{
		"username":    username,
		"nickname":    nickname,
		"password":    model.EncodePassword(password),
		"create_time": now,
		"update_time": now,
	}
	_, err = model.User.Insert(data)
	this.StopE(err, func() {
		this._JSONFail(err.Error())
	}, func() {
		this._JSONSuccess("", nil, this.url("/user/list"))
	})
}

func (this *User) Edit() {
	userID := this.Ctx.GET("id")
	if userID == "" {
		this.Ctx.Redirect("/user/list")
	}
	user, err := model.User.GetByID(userID)
	this.StopE(err, func() {
		this._JSONFail(err.Error())
	})
	this.View.Set("user", user)
	this.View.Layout("admin/form").Render("admin/user/form")
}

func (this *User) Save() {
	nickname := this.Ctx.POST("nickname")
	password := this.Ctx.POST("password")
	password1 := this.Ctx.POST("password1")
	userID := this.Ctx.GET("id")
	if userID == "" {
		this.Ctx.Redirect("/user/list")
	}
	if password != "" {
		if password != password1 {
			this._JSONFail("信息不完整")
		}
		password = model.EncodePassword(password)
	}
	dbUser, err := model.User.GetByID(userID)
	if err != nil || len(dbUser) == 0 {
		this._JSONFail("用户不存在")
	}
	data := gmap.M{
		"nickname":    nickname,
		"update_time": time.Now().Unix(),
	}
	if password != "" {
		data["data"] = password
	}
	_, err = model.User.UpdateByIDs([]string{userID}, data)
	this.StopE(err, func() {
		this._JSONFail(err.Error())
	}, func() {
		this._JSONSuccess("", nil, this.url("/user/list"))
	})
}

func (this *User) Delete() {
	var ids []string
	this.Request.ParseForm()
	id := this.Request.Form["ids"]
	if len(id) > 0 {
		ids = append(ids, id...)
	}
	for _, v := range ids {
		if v == "1" {
			this._JSONFail("系统账号禁止删除")
		}
	}
	_, err := model.User.UpdateByIDs(ids, gmap.M{"is_delete": 1})
	this.StopE(err, func() {
		this._JSONFail(err.Error())
	})
	this._JSONSuccess("", nil, this.url("/user/list"))
}

func (this *User) List() {
	search_field := this.Ctx.GET("search_field")
	keyword := this.Ctx.GET("keyword")
	where := gmap.M{"is_delete": 0}
	rule := map[string]bool{"username": true, "nickname": true, "user_id": true}
	if search_field != "" && keyword != "" && rule[search_field] {
		if search_field == "user_id" {
			where[search_field] = keyword
		} else {
			where[search_field+" like"] = "%" + keyword + "%"
		}
	}
	perPage := gcast.ToInt(this.Ctx.GET("count"))
	if perPage <= 0 || perPage > 100 {
		perPage = 10
	}
	offset := perPage * (gcast.ToInt(this.Ctx.GET("page")) - 1)
	if offset < 0 {
		offset = 0
	}
	users, total, err := model.User.Page(where, offset, perPage)
	this.StopE(err, func() {
		this.WriteE(err)
	})
	this.View.Set("users", users)
	this.View.Set("paginator", this.Ctx.NewPager(perPage, int64(total)))
	this.View.Layout("admin/list").Render("admin/user/list")
}
