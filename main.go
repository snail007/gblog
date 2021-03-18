// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package main

import (
	"fmt"
	emconf "gblog/conf"
	"gblog/initialize"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/http/server"
	gfile "github.com/snail007/gmc/util/file"
	gdaemon "github.com/snail007/gmc/util/process/daemon"
	ghook "github.com/snail007/gmc/util/process/hook"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
)

func main() {

	if err := gdaemon.Start(); err != nil {
		fmt.Println(err)
		return
	}
	if gdaemon.CanRun() {
		//call actual main()
		go start()
	}
	ghook.RegistShutdown(func() {
		gdaemon.Clean()
	})
	ghook.WaitShutdown()
}

func start() {
	var app gcore.App
	defer func() {
		if app != nil && app.Logger() != nil {
			app.Logger().WaitAsyncDone()
		}
	}()
	isDebug := pflag.BoolP("debug", "d", false, "enable debug mode")
	conf := pflag.StringP("conf", "c", "conf/app.toml", "path of config file")
	pflag.Parse()

	if !gfile.Exists(*conf) {
		os.MkdirAll(filepath.Dir(*conf), 0755)
		err := ioutil.WriteFile(*conf, emconf.ConfAPP, 0755)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}
	// 1. create an default app to run.
	app = gmc.New.AppDefault()
	app.Ctx().Set("debug", *isDebug)
	app.SetConfigFile(*conf)
	httpServer := ghttpserver.NewHTTPServer(app.Ctx())
	// 2. add a http server service to app.
	app.AddService(gcore.ServiceItem{
		Service: httpServer,
		AfterInit: func(s *gcore.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = initialize.Initialize(httpServer)
			return
		},
	})

	// 3. run the app
	if err := app.Run(); err != nil {
		fmt.Println(err)
		debug.PrintStack()
	}
}
