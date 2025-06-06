package admin

import (
	"fmt"
	"time"

	"gblog/global"
	"github.com/gookit/validate"
	"github.com/snail007/gmc"
	gcast "github.com/snail007/gmc/util/cast"
	gmap "github.com/snail007/gmc/util/map"
)

type Article struct {
	Admin
}

func (this *Article) List() {
	enableSearch := true
	where := gmap.M{}
	search := this.Ctx.GET("search_field")
	keyword := this.Ctx.GET("keyword")
	catalogID := this.Ctx.GET("catalog_id")
	if catalogID != "" {
		where["catalog_id"] = catalogID
	}
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
		where[search+" like"] = "%" + keyword + "%"
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
	rows, total, err := table.Page(where, start, pageSize, "article_id", "desc")
	if err != nil {
		this.Stop(err)
	}
	catalogTable := gmc.DB.Table("catalog")
	catalogs, err := catalogTable.GetAll("sequence", "asc")
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
		isPrePublish := "0"
		if time.Now().Unix() < gcast.ToInt64(v["create_time"]) {
			isPrePublish = "1"
		}
		rows[k]["is_pre_publish"] = isPrePublish
		rows[k]["catalog_name"] = catalog_name
	}
	this.View.Set("catalogs", catalogs)
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
	isPrePublish := "0"
	if time.Now().Unix() < gcast.ToInt64(row["create_time"]) {
		isPrePublish = "1"
	}
	row["is_pre_publish"] = isPrePublish
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

		v.StringRule("is_draft", "required|enum:0,1")
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
		dataInsert["is_draft"], _ = data.Get("is_draft")
		dataInsert["poster_url"], _ = data.Get("poster_url")
		dataInsert["update_time"] = 0
		t, err := time.ParseInLocation("2006-01-02 15:04:05", this.Ctx.POST("create_time"), time.Local)
		if err != nil {
			t = time.Now()
		} else if t.Before(time.Now()) {
			t = time.Now()
		}
		dataInsert["create_time"] = t.Unix()
		if dataInsert["poster_url"] == "" {
			dataInsert["poster_url"] = global.RandImgIdx()
		}
		id, err := table.Insert(dataInsert)
		if err != nil { // validate ok
			this._JSONFail(err.Error())
		}
		global.Context.Cache().Clear()
		// insert index data
		if global.Context.Indexer() != nil {
			doc := fmt.Sprintf("%s\n%s\n%s", dataInsert["title"], dataInsert["summary"], dataInsert["content"])
			err = global.Context.Indexer().Index(fmt.Sprintf("%d", id), doc)
			if err != nil {
				this.Logger.Warnf("insert index data fail, %d , error: %s", id, err)
			}
		}
		go syncArticle(gcast.ToString(id), nil)
		this._JSONSuccess("", "", this.Ctx.POST("referer"))
	} else {
		catalogTable := gmc.DB.Table("catalog")
		catalogs, err := catalogTable.GetAll("sequence", "asc")
		if err != nil {
			this.Stop(err)
		}
		// show create page
		this.View.Set("showCreateTime", "1")
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
	article, err := table.GetByID(id)
	if err != nil {
		this._JSONFail(err.Error())
	}
	if len(article) == 0 {
		this._JSONFail("data not found")
	}
	if this.Ctx.IsPOST() {
		// do create
		v := data.Create()
		v.FilterRule("article_id", "int")
		v.FilterRule("catalog_id", "int")

		v.StringRule("is_draft", "required|enum:0,1")
		v.StringRule("article_id", "required|min:1")
		v.StringRule("title", "required")
		v.StringRule("summary", "required")
		v.StringRule("content", "required")
		if !v.Validate() { // validate ok
			this._JSONFail(v.Errors.One())
		}

		nowUnix := time.Now().Unix()
		dataUpdate := gmap.M{}
		dataUpdate["title"], _ = data.Get("title")
		dataUpdate["summary"], _ = data.Get("summary")
		dataUpdate["content"], _ = data.Get("content")
		dataUpdate["catalog_id"], _ = data.Get("catalog_id")
		dataUpdate["poster_url"], _ = data.Get("poster_url")
		dataUpdate["update_time"] = nowUnix
		if dataUpdate["poster_url"] == "" {
			dataUpdate["poster_url"] = global.RandImgIdx()
		}
		if article["is_draft"] == "1" {
			dataUpdate["is_draft"], _ = data.Get("is_draft")
			if dataUpdate["is_draft"] == "0" {
				dataUpdate["create_time"] = nowUnix
				dataUpdate["update_time"] = nowUnix
			}
		}
		createTimeDB := gcast.ToInt64(article["create_time"])
		createTime := this.Ctx.POST("create_time")
		now := time.Now()
		if createTime != "" && createTimeDB > now.Unix() {
			t, err := time.ParseInLocation("2006-01-02 15:04:05", createTime, time.Local)
			if err != nil {
				t = now
			} else if t.Before(now) {
				t = now
			}
			dataUpdate["create_time"] = t.Unix()
		}
		_, err = table.UpdateBy(gmap.M{"article_id": id}, dataUpdate)
		if err != nil { // validate ok
			this._JSONFail(err.Error())
		}
		global.Context.Cache().Clear()

		if global.Context.Indexer() != nil {
			// delete & insert index data
			err = global.Context.Indexer().Delete(id)
			if err != nil {
				this.Logger.Warnf("delete index data fail, %d , error: %s", id, err)
			}
			doc := fmt.Sprintf("%s\n%s\n%s", dataUpdate["title"], dataUpdate["summary"], dataUpdate["content"])
			err = global.Context.Indexer().Index(id, doc)
			if err != nil {
				this.Logger.Warnf("insert index data fail, %d , error: %s", id, err)
			}
		}
		go syncArticle(id, article)
		this._JSONSuccess("", "", this.Ctx.POST("referer"))
	} else {
		catalogTable := gmc.DB.Table("catalog")
		catalogs, err := catalogTable.GetAll("sequence", "asc")
		if err != nil {
			this.Stop(err)
		}
		// show page
		createTime := gcast.ToInt64(article["create_time"])
		showCreateTime := "0"
		if createTime > time.Now().Unix() {
			showCreateTime = "1"
		}
		this.View.Set("showCreateTime", showCreateTime)
		this.View.Set("catalogs", catalogs)
		this.View.Set("data", article)
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
	var articles []gmap.Mss
	var err error
	if isGitLabSync() {
		articles, err = table.MGetByIDs(ids)
	}
	_, err = table.DeleteByIDs(ids)
	this.StopE(err, func() {
		this._JSONFail(err.Error())
	})
	global.Context.Cache().Clear()
	if global.Context.Indexer() != nil {
		for _, id := range ids {
			//delete index data
			err = global.Context.Indexer().Delete(id)
			if err != nil {
				this.Logger.Warnf("delete index data fail, %d , error: %s", id, err)
			}
		}
	}
	if len(articles) > 0 {
		go func() {
			for _, v := range articles {
				syncDeleteArticle(v)
			}
		}()
	}
	this._JSONSuccess("", nil, this.Ctx.Header("Referer"))
}

func (this *Article) Move() {
	var ids []string
	this.Request.ParseForm()
	id := this.Request.Form["ids"]
	if len(id) > 0 {
		ids = append(ids, id...)
	}
	catalogID := this.Ctx.POST("catalog_id")
	catalogTable := gmc.DB.Table("catalog")
	catalog, err := catalogTable.GetByID(catalogID)
	if err != nil {
		this._JSONFail(err.Error())
	}
	if len(catalog) == 0 {
		this._JSONFail("catalog not found")
	}
	table := gmc.DB.Table("article")
	var articles []gmap.Mss
	if isGitLabSync() {
		articles, err = table.MGetByIDs(ids)
	}
	_, err = table.UpdateByIDs(ids, gmap.M{"catalog_id": catalogID})
	this.StopE(err, func() {
		this._JSONFail(err.Error())
	})
	if len(articles) > 0 {
		go func() {
			for _, v := range articles {
				syncDeleteArticle(v)
			}
			articlesNew, _ := table.MGetByIDs(ids)
			for _, v := range articlesNew {
				syncArticle(v["article_id"], nil)
			}
		}()
	}
	global.Context.Cache().Clear()
	if isGitLabSync() {

	}
	this._JSONSuccess("", nil, this.Ctx.Header("Referer"))
}
