package admin

import (
	"github.com/snail007/gmc"
	gmccaptcha "github.com/snail007/gmc/util/captcha"
	gmap "github.com/snail007/gmc/util/map"
	"image/png"
	"gblog/model"
	"strings"
)

var (
	cap    = gmc.New.CaptchaDefault()
	capKEY = "admin_captcha"
)

type Login struct {
	Base
}

func (this *Login) Auth() {
	err := this.SessionStart()
	if err != nil {
		this.Stop(err)
	}
	u := this.Session.Get("admin")
	if u != nil && u.(gmap.Mss)["username"] != "" {
		this._JSONSuccess("", "", this.url("/main/index"))
	}
	captcha := strings.TrimSpace(this.Ctx.POST("captcha"))
	captchaSession0 := this.Session.Get(capKEY)
	captchaSession := ""
	if captchaSession0 != nil {
		captchaSession = captchaSession0.(string)
	}

	this.Session.Delete(capKEY)
	if captchaSession == "" || captcha == "" || captchaSession != strings.ToLower(captcha) {
		this._JSONFail("验证码错误")
	}

	username := this.Ctx.POST("username")
	password := this.Ctx.POST("password")
	if password == "" || username == "" {
		this._JSONFail("信息不完整")
	}
	dbUser, err := model.User.GetBy(gmap.M{"username": username})
	if err != nil || len(dbUser) == 0 {
		this._JSONFail("用户名或密码错误")
	}
	if dbUser["password"] != model.EncodePassword(password) {
		this._JSONFail("用户名或密码错误")
	}
	delete(dbUser, "password")
	this.Session.Set("admin", dbUser)
	this._JSONSuccess("", "", this.url("/main/index"))
}

func (this *Login) Index_() {
	this.View.Set("admPath",this.admPath())
	this.View.Layout("admin/login").Render("admin/login/login")
}

func (this *Login) Logout() {
	this.SessionStart()
	this.SessionDestroy()
	this.Ctx.Redirect(this.url("/"))
}

func (this *Login) Captcha() {
	this.SessionStart()
	img, str := cap.Create(4, gmccaptcha.CLEAR)
	this.Session.Set(capKEY, strings.ToLower(str))
	png.Encode(this.Response, img)
}
