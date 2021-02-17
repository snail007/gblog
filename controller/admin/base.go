package admin

import (
	"github.com/snail007/gmc"
	gmap "github.com/snail007/gmc/util/map"
)

type Base struct {
	gmc.Controller
}

func (this *Base) url(url string) string {
	return this.admPath() + url
}

func (this *Base) admPath() string {
	return this.Config.GetString("admin.urlpath")
}

func (this *Base) _JSONSuccess(msg string, data ...interface{}) {
	this._JSON(msg, 200, data...)
}

func (this *Base) _JSONFail(msg string) {
	this._JSON(msg, 500)
}

func (this *Base) _JSON(msg string, code int, data ...interface{}) {
	var data0 interface{}
	var url1 interface{}
	if len(data) >= 1 {
		data0 = data[0]
	}
	if len(data) >= 2 {
		url1 = data[1]
	}
	d := gmap.M{
		"msg":  msg,
		"code": code,
		"data": data0,
		"url":  url1,
	}
	this.Write(d)
	this.Stop()
}
