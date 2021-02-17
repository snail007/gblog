package admin

import (
	gmap "github.com/snail007/gmc/util/map"
)

type Admin struct {
	Base
	User gmap.Mss
}

func (this *Admin) Before() {
	err := this.SessionStart()
	if err != nil {
		this.Stop(err)
	}
	if u, ok := this._IsLogin(); !ok {
		this.Ctx.Redirect(this.url("/"))
		return
	} else {
		this.User = u
	}
	this.View.Set("admPath",this.admPath())
	this.View.Set("admin",this.User)
}

func (this *Admin) _IsLogin() (user gmap.Mss, isLogin bool) {
	u := this.Session.Get("admin")
	if u != nil && u.(gmap.Mss)["username"] != "" {
		return u.(gmap.Mss), true
	}
	return nil, false
}
