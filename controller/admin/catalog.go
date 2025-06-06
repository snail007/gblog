package admin

import (
	"gblog/global"
	"github.com/gookit/validate"
	"github.com/snail007/gmc"
	gdb "github.com/snail007/gmc/module/db"
	gmap "github.com/snail007/gmc/util/map"
)

type Catalog struct {
	Admin
}

func (this *Catalog) List() {
	enableSearch := false
	where := gmap.M{}
	search := this.Ctx.GET("search_field")
	keyword := this.Ctx.GET("keyword")
	if enableSearch && search != "" && keyword != "" {
		data, err := validate.FromRequest(this.Request)
		if err != nil {
			this.Stop(err)
		}
		v := data.Create()
		v.StringRule("search_field", "enum:name")
		if !v.Validate() {
			this.Stop(v.Errors.One())
		}
		where = gmap.M{search + " like": "%" + keyword + "%"}
	}
	table := gmc.DB.Table("catalog")
	rows, err := table.MGetBy(where, "is_nav", "desc", "sequence", "asc")
	if err != nil {
		this.Stop(err)
	}
	db := gdb.DB()
	rs, err := db.Query(db.AR().Select("count(*) as total,catalog_id").From("article").GroupBy("catalog_id"))
	if err != nil {
		this.Stop(err)
	}
	catalogsSummary := rs.MapRows("catalog_id")
	for k, v := range rows {
		cnt := ""
		vv, exists := catalogsSummary[v["catalog_id"]]
		if !exists {
			cnt = "0"
		} else {
			cnt = vv["total"]
		}
		v["total"] = cnt
		rows[k] = v
	}
	this.View.Set("rows", rows)
	this.View.Set("enable_search", enableSearch)
	this.View.Layout("admin/list").Render("admin/catalog/list")
}

func (this *Catalog) Detail() {
	id := this.Ctx.GET("catalog_id")
	if id == "" {
		id = this.Ctx.POST("catalog_id")
	}
	table := gmc.DB.Table("catalog")
	row, err := table.GetByID(id)
	if err != nil {
		this._JSONFail(err.Error())
	}
	if len(row) == 0 {
		this._JSONFail("data not found")
	}
	this.View.Set("data", row)
	this.View.Layout("admin/list").Render("admin/catalog/detail")
}

func (this *Catalog) Create() {
	if this.Ctx.IsPOST() {
		// do create
		data, err := validate.FromRequest(this.Request)
		if err != nil {
			this.Stop(err)
		}
		v := data.Create()
		v.FilterRule("sequence", "int")

		v.StringRule("name", "required|minLen:1")
		v.StringRule("sequence", "required|minLen:1|min:0")
		v.StringRule("is_nav", "required|enum:0,1")
		if !v.Validate() { // validate ok
			this._JSONFail(v.Errors.One())
		}
		table := gmc.DB.Table("catalog")
		dataInsert := gmap.M{}
		dataInsert["name"], _ = data.Get("name")
		dataInsert["sequence"], _ = data.Get("sequence")
		dataInsert["is_nav"], _ = data.Get("is_nav")
		_, err = table.Insert(dataInsert)
		if err != nil { // validate ok
			this._JSONFail(err.Error())
		}
		global.Context.Cache().Clear()
		this._JSONSuccess("", "", this.Ctx.POST("referer"))
	} else {
		// show create page
		this.View.Layout("admin/form").Render("admin/catalog/form")
	}
}

func (this *Catalog) Edit() {
	data, err := validate.FromRequest(this.Request)
	if err != nil {
		this.Stop(err)
	}
	table := gmc.DB.Table("catalog")
	id := this.Ctx.GET("catalog_id")
	if id == "" {
		id = this.Ctx.POST("catalog_id")
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
		v.FilterRule("catalog_id", "int")
		v.FilterRule("sequence", "int")

		v.StringRule("is_nav", "required|enum:0,1")
		v.StringRule("name", "required|minLen:1")
		if !v.Validate() { // validate ok
			this._JSONFail(v.Errors.One())
		}

		dataInsert := gmap.M{}
		dataInsert["name"], _ = data.Get("name")
		dataInsert["sequence"], _ = data.Get("sequence")
		dataInsert["is_nav"], _ = data.Get("is_nav")
		_, err = table.UpdateBy(gmap.M{"catalog_id": id}, dataInsert)
		if err != nil { // validate ok
			this._JSONFail(err.Error())
		}
		global.Context.Cache().Clear()
		this._JSONSuccess("", "", this.Ctx.POST("referer"))
	} else {
		// show create page
		this.View.Set("data", row)
		this.View.Layout("admin/form").Render("admin/catalog/form")
	}
}

// 保存拖拽排序
func (this *Catalog) SaveOrder() {
	data := this.Ctx.POSTData()
	if len(data) == 0 {
		this._JSONFail("数据不能为空！")
	}
	d := []gmap.M{}
	db := gmc.DB.DB()
	for k, v := range data {
		d = append(d, map[string]interface{}{
			"catalog_id": k,
			"sequence":   v,
		})
	}
	_, err := db.Exec(db.AR().UpdateBatch("catalog", d, []string{"catalog_id"}))
	if err != nil {
		this._JSONFail("修改排序失败！" + err.Error())
	} else {
		this._JSONSuccess("")
	}
}

func (this *Catalog) Delete() {
	var ids []string
	this.Request.ParseForm()
	id := this.Request.Form["ids"]
	if len(id) > 0 {
		ids = append(ids, id...)
	}
	for _, v := range ids {
		if v == "0" {
			this._JSONFail("默认分类不能删除")
		}
	}
	table := gmc.DB.Table("catalog")
	_, err := table.DeleteByIDs(ids)
	this.StopE(err, func() {
		this._JSONFail(err.Error())
	})
	gmc.DB.Table("article").UpdateBy(gmap.M{"catalog_id": ids}, gmap.M{"catalog_id": 0})
	global.Context.Cache().Clear()
	this._JSONSuccess("", nil, this.Ctx.Header("Referer"))
}
