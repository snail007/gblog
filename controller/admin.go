// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package controller

import (
	"github.com/snail007/gmc"
)

type Admin struct {
	gmc.Controller
}

func (this *Admin) Before() {
	if this.Ctx.ControllerMethod() == "Login" {
		return
	}
	err := this.SessionStart()
	if err != nil {
		this.Die(err)
	}
	if !this.isLogin() {
		this.Ctx.Redirect("/login")
	}
}

func (this *Admin) Login() {
	this.Write("login")
}

func (this *Admin) isLogin() bool {
	return false
}
