// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package main

import (
	"fmt"
	"gblog/initialize"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/http/server"
	"runtime/debug"
)

func main() {

	// 1. create an default app to run.
	app := gmc.New.AppDefault()
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
