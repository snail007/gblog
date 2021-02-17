package admin

import (
	"gblog/global"
	"github.com/gookit/validate"
	"github.com/snail007/gmc"
	gcast "github.com/snail007/gmc/util/cast"
	gmap "github.com/snail007/gmc/util/map"
	"time"
)

type Article struct {
	Admin
}

func (this *Article) List() {
	enableSearch := true
	where := gmap.M{}
	search := this.Ctx.GET("search_field")
	keyword := this.Ctx.GET("keyword")
	if enableSearch && search != "" && keyword != "" {
		data, err := validate.FromRequest(this.Request)
		if err != nil {
			this.Stop(err)
		}
		v := data.Create()
		v.StringRule("search_field", "enum:title,content")
		if !v.Validate() {
			this.Stop(v.Errors.One())
		}
		where = gmap.M{search + " like": "%" + keyword + "%"}
	}
	page := gcast.ToInt(this.Ctx.GET("page"))
	pageSize := gcast.ToInt(this.Ctx.GET("count"))
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	start := (page - 1) * pageSize
	if start < 0 {
		start = 0
	}
	table := gmc.DB.Table("article")
	rows, total, err := table.Page(where, start, pageSize, gmap.M{"article_id": "desc"})
	if err != nil {
		this.Stop(err)
	}
	catalogTable := gmc.DB.Table("catalog")
	catalogs, err := catalogTable.GetAll(gmap.M{"sequence": "asc"})
	if err != nil {
		this.Stop(err)
	}
	for k, v := range rows {
		catalog_name := "未知分类"
		for _, c := range catalogs {
			if v["catalog_id"] == c["catalog_id"] {
				catalog_name = c["name"]
			}
		}
		rows[k]["catalog_name"] = catalog_name
	}
	this.View.Set("rows", rows)
	this.View.Set("enable_search", enableSearch)
	this.View.Set("paginator", this.Ctx.NewPager(pageSize, int64(total)))
	this.View.Layout("admin/list").Render("admin/article/list")
}

func (this *Article) Detail() {
	id := this.Ctx.GET("article_id")
	if id == "" {
		id = this.Ctx.POST("article_id")
	}
	table := gmc.DB.Table("article")
	row, err := table.GetByID(id)
	if err != nil {
		this._JSONFail(err.Error())
	}
	if len(row) == 0 {
		this._JSONFail("data not found")
	}
	catalogTable := gmc.DB.Table("catalog")
	catalog, err := catalogTable.GetByID(row["catalog_id"])
	if err != nil {
		this.Stop(err)
	}
	row["catalog_name"] = catalog["name"]
	this.View.Set("data", row)
	this.View.Layout("admin/form").Render("admin/article/detail")
}

func (this *Article) Create() {
	if this.Ctx.IsPOST() {
		// do create
		data, err := validate.FromRequest(this.Request)
		if err != nil {
			this.Stop(err)
		}
		v := data.Create()
		v.FilterRule("catalog_id", "int")
		v.StringRule("content", "required")
		v.StringRule("title", "required")
		v.StringRule("summary", "required")
		if !v.Validate() { // validate ok
			this._JSONFail(v.Errors.One())
		}
		table := gmc.DB.Table("article")
		dataInsert := gmap.M{}
		dataInsert["title"], _ = data.Get("title")
		dataInsert["summary"], _ = data.Get("summary")
		dataInsert["content"], _ = data.Get("content")
		dataInsert["catalog_id"], _ = data.Get("catalog_id")
		dataInsert["poster_url"], _ = data.Get("poster_url")
		dataInsert["update_time"] = 0
		dataInsert["create_time"] = time.Now().Unix()
		if dataInsert["poster_url"] == "" {
			dataInsert["poster_url"] = global.RandImgIdx()
		}
		_, err = table.Insert(dataInsert)
		if err != nil { // validate ok
			this._JSONFail(err.Error())
		}
		this._JSONSuccess("", "", this.url("/article/list"))
	} else {
		catalogTable := gmc.DB.Table("catalog")
		catalogs, err := catalogTable.GetAll(gmap.M{"sequence": "asc"})
		if err != nil {
			this.Stop(err)
		}
		// show create page
		this.View.Set("catalogs", catalogs)
		this.View.Layout("admin/form").Render("admin/article/form")
	}
}

func (this *Article) Edit() {
	data, err := validate.FromRequest(this.Request)
	if err != nil {
		this.Stop(err)
	}
	table := gmc.DB.Table("article")
	id := this.Ctx.GET("article_id")
	if id == "" {
		id = this.Ctx.POST("article_id")
	}
	row, err := table.GetByID(id)
	if err != nil {
		this._JSONFail(err.Error())
	}
	if len(row) == 0 {
		this._JSONFail("data not found")
	}
	if this.Ctx.IsPOST() {
		// do create
		v := data.Create()
		v.FilterRule("article_id", "int")
		v.FilterRule("catalog_id", "int")

		v.StringRule("article_id", "required|min:1")
		v.StringRule("title", "required")
		v.StringRule("summary", "required")
		v.StringRule("content", "required")
		if !v.Validate() { // validate ok
			this._JSONFail(v.Errors.One())
		}

		dataUpdate := gmap.M{}
		dataUpdate["title"], _ = data.Get("title")
		dataUpdate["summary"], _ = data.Get("summary")
		dataUpdate["content"], _ = data.Get("content")
		dataUpdate["catalog_id"], _ = data.Get("catalog_id")
		dataUpdate["poster_url"], _ = data.Get("poster_url")
		dataUpdate["update_time"] = time.Now().Unix()
		if dataUpdate["poster_url"] == "" {
			dataUpdate["poster_url"] = global.RandImgIdx()
		}
		_, err = table.UpdateBy(gmap.M{"article_id": id}, dataUpdate)
		if err != nil { // validate ok
			this._JSONFail(err.Error())
		}
		this._JSONSuccess("", "", this.url("/article/list"))
	} else {
		catalogTable := gmc.DB.Table("catalog")
		catalogs, err := catalogTable.GetAll(gmap.M{"sequence": "asc"})
		if err != nil {
			this.Stop(err)
		}
		// show create page
		this.View.Set("catalogs", catalogs)
		this.View.Set("data", row)
		this.View.Layout("admin/form").Render("admin/article/form")
	}
}

func (this *Article) Delete() {
	var ids []string
	this.Request.ParseForm()
	id := this.Request.Form["ids"]
	if len(id) > 0 {
		ids = append(ids, id...)
	}
	table := gmc.DB.Table("article")
	_, err := table.DeleteByIDs(ids)
	this.StopE(err, func() {
		this._JSONFail(err.Error())
	})
	this._JSONSuccess("", nil, this.url("/article/list"))
}
