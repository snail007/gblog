// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package initialize

import (
	"gblog/global"
	"gblog/router"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gcache "github.com/snail007/gmc/module/cache"
	gdb "github.com/snail007/gmc/module/db"
	glog "github.com/snail007/gmc/module/log"
	gfile "github.com/snail007/gmc/util/file"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

func Initialize(s gcore.HTTPServer) (err error) {
	defer func() {
		if err != nil {
			err = gmc.Err.Wrap(err)
		}
	}()
	// init command line
	isDebug := pflag.BoolP("debug", "d", false, "enable debug mode")
	conf := pflag.StringP("conf", "c", "conf/app.toml", "path of config file")
	pflag.Parse()

	// init Context
	ctx, err := global.NewBContext(*conf)
	if err != nil {
		return
	}
	ctx.SetIsDebug(*isDebug)
	ctx.SetServer(s)
	cfg := ctx.Config()

	//init logger
	ctx.SetLog(s.Logger())

	// init db directory
	dataFile := cfg.Get("database.sqlite3").([]interface{})[0].(map[string]interface{})["database"].(string)
	dataDir := filepath.Dir(dataFile)
	if !gfile.Exists(dataDir) {
		err = os.MkdirAll(dataDir, 0700)
		if err != nil {
			return
		}
	}

	// init db
	err = gdb.Init(cfg)
	if err != nil {
		return
	}
	ctx.SetDB(gdb.DB())

	// auto init databases
	err = checkTable(ctx.DB())
	if err != nil {
		return
	}

	// init logger
	logger := glog.NewFromConfig(cfg, "")
	if ctx.IsDebug() {
		logger.SetLevel(gcore.LTRACE)
	}
	ctx.SetLog(logger)

	//init cache
	err = gcache.Init(cfg)
	if err != nil {
		return
	}
	ctx.SetCache(gcache.Cache())

	global.Context = ctx

	// init router
	router.InitRouter(s)
	return
}

func checkTable(db gcore.Database) (err error) {
	_, e := db.Query(db.AR().Raw("select * from article"))
	if e == nil {
		return
	}
	// create table
	_, err = db.Exec(db.AR().Raw("create table article(id integer PRIMARY KEY AUTOINCREMENT,title text,content text,catalog_id int,create_time int,update_time int)"))
	if err != nil {
		return
	}
	_, err = db.Exec(db.AR().Raw("create table catalog(id integer PRIMARY KEY AUTOINCREMENT,name text)"))
	if err != nil {
		return
	}
	return
}
