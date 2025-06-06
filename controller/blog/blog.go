package blog

import (
	"encoding/json"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gblog/global"
	"github.com/blevesearch/bleve"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gcast "github.com/snail007/gmc/util/cast"
	gfile "github.com/snail007/gmc/util/file"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/xanzy/go-gitlab"
)

type Blog struct {
	gmc.Controller
	bConf map[string]gmap.M
}

var (
	navMethod = map[string]bool{
		"List":     true,
		"Views":    true,
		"Timeline": true,
		"Search":   true,
		"Catalogs": true,
	}
	cacheSecond uint = 3600
)

func (this *Blog) Before() {
	this.bConf = map[string]gmap.M{}

	//init config
	db := gmc.DB.DB()
	rs, err := db.Query(this.cache("allConfig").From("config"))
	if err != nil {
		this.Stop(err)
	}
	for _, v := range rs.Rows() {
		value := gmap.M{}
		err := json.Unmarshal([]byte(v["value"]), &value)
		if err != nil {
			value = gmap.M{}
		}
		this.bConf[v["key"]] = value
		this.View.Set("bc_"+v["key"], value)
	}
	//maintain checking
	status := gcast.ToString(this.bConf["basic"]["web_site_status"])
	if status != "on" {
		this.View.Render("blog/maintain")
		this.Die()
	}
	// init nav
	method := this.Ctx.ControllerMethod()
	if navMethod[method] {
		this.buildNav()
	}
}

func (this *Blog) cache(key string) (ar gcore.ActiveRecord) {
	ar = gmc.DB.DB().AR()
	if cacheSecond <= 0 {
		return ar
	}
	ar.Cache(key, cacheSecond+uint(rand.Int31n(30)))
	return ar
}
func (this *Blog) buildNav() {
	db := gmc.DB.DB()
	rs, err := db.Query(this.cache("nav").From("catalog").Where(gmap.M{"is_nav": 1}).OrderBy("sequence", "asc"))
	if err != nil {
		this.Stop(err)
	}
	navs := rs.Rows()
	this.View.Set("navs", navs)
}

func (this *Blog) List() {
	page := gcast.ToInt(this.Ctx.GET("page"))
	id := this.Param.ByName("id")
	where := gmap.M{
		"create_time <=": time.Now().Unix(),
		"is_draft":       0,
	}
	db := gmc.DB.DB()
	catalog := gmap.Mss{}
	if id != "" {
		where["catalog_id"] = id
		rs, err := db.Query(this.cache("catalog-" + id).From("catalog").Where(gmap.M{"catalog_id": id}))
		if err != nil {
			this.Stop(err)
		}
		catalog = rs.Row()
	}
	if id != "" && len(catalog) == 0 {
		this.Ctx.WriteHeader(404)
		return
	}
	rs, err := db.Query(this.cache("articles-"+id).From("article").Where(where).OrderBy("create_time", "desc"))
	if err != nil {
		this.Stop(err)
	}
	articles := rs.Rows()

	rs, err = db.Query(this.cache("allCatalog").From("catalog"))
	if err != nil {
		this.Stop(err)
	}
	catalogs := rs.MapRows("catalog_id")
	for k, v := range articles {
		catalog := catalogs[v["catalog_id"]]
		v["catalog_name"] = catalog["name"]
		articles[k] = v
	}
	// page
	if page <= 0 {
		page = 1
	}
	pageSize := 20
	start := (page - 1) * pageSize
	end := pageSize * page
	total := len(articles)
	if start < 0 {
		start = 0
	}
	if end > total {
		end = total
	}
	if start > end {
		start = 0
		end = pageSize
	}
	if pageSize > total {
		pageSize = total
	}
	pArticles := articles[start:end]
	this.View.Set("title", catalog["name"])
	this.View.Set("catalog", catalog)
	this.View.Set("articles", pArticles)
	this.View.Set("paginator", this.Ctx.NewPager(pageSize, int64(total)))
	this.View.Layout("blog/list").Render("blog/list")
}

func (this *Blog) Views() {
	id := this.Param.ByName("id")
	db := gmc.DB.DB()
	rs, err := db.Query(this.cache("article-" + id).From("article").Where(
		gmap.M{"article_id": id, "create_time <=": time.Now().Unix(), "is_draft": 0},
	))
	if err != nil {
		this.Stop(err)
	}
	article := rs.Row()

	if len(article) == 0 {
		this.Ctx.WriteHeader(404)
		return
	}

	rs, err = db.Query(this.cache("catalog-" + id).From("catalog").Where(gmap.M{"catalog_id": article["catalog_id"]}))
	if err != nil {
		this.Stop(err)
	}
	catalog := rs.Row()
	if len(catalog) == 0 {
		this.Ctx.WriteHeader(404)
		return
	}

	articlePre, articleNext := gmap.Mss{}, gmap.Mss{}

	rs, err = db.Query(this.cache("articlePre-"+article["article_id"]).
		From("article").
		Where(gmap.M{"create_time <": article["create_time"], "catalog_id": article["catalog_id"]}).
		OrderBy("create_time", "desc").Limit(1))
	if err != nil {
		this.Stop(err)
	}
	articlePre = rs.Row()

	rs, err = db.Query(this.cache("articlesNext-"+article["article_id"]).
		From("article").
		Where(gmap.M{"create_time >": article["create_time"], "catalog_id": article["catalog_id"]}).
		OrderBy("create_time", "asc").Limit(1))
	if err != nil {
		this.Stop(err)
	}
	articleNext = rs.Row()

	this.View.Set("title", article["title"])
	this.View.Set("article", article)
	this.View.Set("articlePre", articlePre)
	this.View.Set("articleNext", articleNext)
	this.View.Set("catalog", catalog)
	this.View.Layout("blog/content").Render("blog/content")
}

func (this *Blog) Timeline() {
	db := gmc.DB.DB()
	rs, err := db.Query(this.cache("timeline").
		From("article").
		Where(gmap.M{"create_time <=": time.Now().Unix(), "is_draft": 0}).
		OrderBy("create_time", "desc"))
	if err != nil {
		this.Stop(err)
	}
	articles := rs.Rows()
	this.View.Set("articles", articles)
	this.View.Layout("blog/timeline").Render("blog/timeline")
}
func (this *Blog) Search() {
	keyword := strings.Trim(this.Ctx.GET("keyword"), " \t")
	if keyword == "" || len(keyword) == 1 || len(keyword) >= 100 {
		this.Ctx.WriteHeader(http.StatusNotFound)
		return
	}
	var articles = []gmap.Mss{}

	// sql search
	db := gmc.DB.DB()
	rs, err := db.Query(this.cache("search").
		From("article").
		Where(gmap.M{"create_time <=": time.Now().Unix(), "is_draft": 0}).
		OrderBy("create_time", "desc"))
	if err != nil {
		this.Stop(err)
	}
	articlesAll := rs.Rows()
	titleMatch, summaryMatch, contentMatch := []gmap.Mss{}, []gmap.Mss{}, []gmap.Mss{}
	keyword = strings.ToLower(keyword)
	for _, v := range articlesAll {
		if strings.Contains(strings.ToLower(v["title"]), keyword) {
			titleMatch = append(titleMatch, v)
		} else if strings.Contains(strings.ToLower(v["summary"]), keyword) {
			summaryMatch = append(summaryMatch, v)
		} else if strings.Contains(strings.ToLower(v["content"]), keyword) {
			contentMatch = append(contentMatch, v)
		}
	}
	articles = append(titleMatch, summaryMatch...)
	articles = append(articles, contentMatch...)

	// sql search empty, try bleve search
	if len(articles) == 0 && global.Context.Indexer() != nil && this.Ctx.Config().GetBool("search.enablefulltextindex") {
		req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(keyword))
		req.Size = 100
		req.Highlight = bleve.NewHighlight()
		res, err := global.Context.Indexer().Search(req)
		if err != nil {
			this.Stop(err)
		}
		if res.Total > 0 {
			if res.Request.Size > 0 {
				articlesIDArr := []string{}
				for _, hit := range res.Hits {
					articlesIDArr = append(articlesIDArr, hit.ID)
				}
				db := gmc.DB.DB()
				rs, err := db.Query(this.cache("search").
					From("article").
					Where(gmap.M{"create_time <=": time.Now().Unix(), "is_draft": 0}).
					OrderBy("create_time", "desc"))
				if err != nil {
					this.Stop(err)
				}
				rows := rs.MapRows("article_id")
				for _, articleID := range articlesIDArr {
					article, ok := rows[articleID]
					if !ok {
						continue
					}
					articles = append(articles, article)
				}
			}
		}
	}
	if len(articles) > 10 {
		articles = articles[:10]
	}
	this.View.Set("articles", articles)
	this.View.Layout("blog/timeline").Render("blog/timeline")
}

func (this *Blog) Catalogs() {
	db := gmc.DB.DB()
	rs, err := db.Query(this.cache("catalogs").From("catalog").OrderBy("sequence", "asc"))
	if err != nil {
		this.Stop(err)
	}
	catalogs := rs.Rows()

	rs, err = db.Query(this.cache("catalogsSummary").
		Select("count(*) as total,catalog_id").
		From("article").
		Where(gmap.M{"create_time <=": time.Now().Unix(), "is_draft": 0}).
		GroupBy("catalog_id"))
	if err != nil {
		this.Stop(err)
	}
	catalogsSummary := rs.MapRows("catalog_id")
	for k, v := range catalogs {
		cnt := ""
		vv, exists := catalogsSummary[v["catalog_id"]]
		if !exists {
			cnt = "0"
		} else {
			cnt = vv["total"]
		}
		v["total"] = cnt
		v["rand"] = global.RandImgIdx()
		catalogs[k] = v
	}
	this.View.Set("catalogs", catalogs)
	this.View.Layout("blog/list").Render("blog/catalogs")
}

func (this *Blog) Attachment() {
	id := this.Ctx.GET("id")
	if id == "" {
		id = this.Ctx.GetParam("id")
	}
	if id == "" {
		this.Ctx.WriteHeader(http.StatusNotFound)
		return
	}
	id = filepath.Clean(strings.TrimPrefix(id, "/"))
	path := "attachment/" + id
	storageType := global.Context.BConfig("upload.upload_file_storage")
	switch storageType {
	case "local":
		rootDir := gfile.Abs(this.Config.GetString("attachment.dir"))
		file := gfile.Abs(filepath.Join(rootDir, filepath.Clean(id)))
		if !strings.Contains(file, rootDir) {
			this.Ctx.WriteHeader(http.StatusNotFound)
			return
		}
		this.Ctx.WriteFile(file)
	case "github":
		userRepo := gcast.ToString(global.Context.BConfig("upload.github_repo"))
		speedURL := gcast.ToString(global.Context.BConfig("upload.github_speed_url"))
		if speedURL == "" {
			speedURL = "https://cdn.jsdelivr.net/gh/%u/%p"
		}
		speedURL = strings.Replace(speedURL, "%u", userRepo, 1)
		speedURL = strings.Replace(speedURL, "%p", path, 1)
		this.Ctx.Redirect(speedURL)
	case "gitlab":
		mimeType := mime.TypeByExtension(filepath.Ext(id))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		var bin []byte
		if gfile.Exists(path) {
			bin = gfile.Bytes(path)
		} else {
			token := gcast.ToString(global.Context.BConfig("upload.gitlab_token"))
			apiURL := gcast.ToString(global.Context.BConfig("upload.gitlab_api_url"))
			userRepo := gcast.ToString(global.Context.BConfig("upload.gitlab_repo"))
			client, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL))
			if err != nil {
				this.Logger.Warnf("create gitlab client error, %s", err)
				this.Ctx.Write(err.Error())
				return
			}
			var resp *gitlab.Response
			bin, resp, err = client.RepositoryFiles.GetRawFile(userRepo, path, &gitlab.GetRawFileOptions{})
			if err != nil {
				this.Logger.Warnf("get gitlab file error, %s", err)
				this.Ctx.Write(err.Error())
				return
			}
			resp.Body.Close()
		}
		this.Ctx.SetHeader("Content-Type", mimeType)
		this.Ctx.Write(bin)
		go func() {
			if !gfile.Exists(path) {
				dir := filepath.Dir(path)
				if !gfile.IsDir(dir) {
					os.MkdirAll(dir, 0755)
				}
				err := gfile.Write(path, bin, false)
				if err != nil {
					this.Logger.Warnf("write gitlab file to local fail, error: %s", err)
				}
			}
		}()
	}
}
