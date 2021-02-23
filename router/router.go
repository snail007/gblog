// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package router

import (
	"gblog/controller/admin"
	"gblog/controller/blog"
	gcore "github.com/snail007/gmc/core"
)

func InitRouter(s gcore.HTTPServer) {
	r := s.Router()
	admPath := s.Config().GetString("admin.urlpath")
	urlPath := s.Config().GetString("attachment.url")
	adm := r.Group(admPath)
	adm.Controller("/admin", new(admin.Admin))
	adm.ControllerMethod("/", new(admin.Login), "Index_")
	adm.Controller("/login", new(admin.Login))
	adm.Controller("/main", new(admin.Main))
	adm.Controller("/user", new(admin.User))
	adm.Controller("/catalog", new(admin.Catalog))
	adm.Controller("/article", new(admin.Article))
	adm.Controller("/config", new(admin.Config))
	adm.Controller("/attachment", new(admin.Attachment))

	r.ControllerMethod("/", new(blog.Blog), "List")
	r.ControllerMethod("/list/:id", new(blog.Blog), "List")
	r.ControllerMethod("/view/:id", new(blog.Blog), "Views")
	r.ControllerMethod("/timeline", new(blog.Blog), "Timeline")
	r.ControllerMethod("/search", new(blog.Blog), "Search")
	r.ControllerMethod("/catalogs", new(blog.Blog), "Catalogs")
	r.ControllerMethod(urlPath, new(blog.Blog), "Attachment")
	r.ControllerMethod(urlPath+"/*id", new(blog.Blog), "Attachment")

	// indicates router initialized
	s.Logger().Infof("router inited.")
}
