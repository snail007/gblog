package admin

import (
	"github.com/snail007/gmc"
	gmap "github.com/snail007/gmc/util/map"
)

type Main struct {
	Admin
}

func (this *Main) Index() {
	this.View.Layout("admin/index").Render("admin/main/index")
}

func (this *Main) Main() {
	db := gmc.DB.DB()
	rs, err := db.Query(db.AR().Select("count(*) as total").From("article"))
	if err != nil {
		this.Stop(err)
	}
	stat := gmap.M{
		"articleCount": rs.Value("total"),
	}
	this.View.Set("stat", stat)
	this.View.Layout("admin/page").Render("admin/main/main")
}

func (this *Main) Other() {
	this.View.Layout("admin/list").Render("admin/other/code")
}
