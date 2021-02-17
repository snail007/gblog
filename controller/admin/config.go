package admin

import (
	"encoding/json"
	"fmt"
	"gblog/global"
	"github.com/snail007/gmc"
	gmap "github.com/snail007/gmc/util/map"
)

type Config struct {
	Admin
}

var (
	keyList = map[string]bool{
		"basic":  true,
		"system": true,
		"upload": true,
	}
)

func (this *Config) configValue(key string) (value gmap.M, err error) {
	value =gmap.M{}
	table := gmc.DB.Table("config")
	row, err := table.GetBy(gmap.M{"key": key})
	if err != nil {
		return
	}
	if len(row) == 0 {
		err = fmt.Errorf("data not found")
		return
	}
	err = json.Unmarshal([]byte(row["value"]), &value)
	return
}
func (this *Config) Conf() {
	key := this.Ctx.GetPost("key", "basic")
	if !keyList[key] {
		this._JSONFail("key not found")
	}
	if this.Ctx.IsPOST() {
		data := this.Ctx.POSTData()
		dStr, err := json.Marshal(data)
		if err != nil {
			this._JSONFail(err.Error())
		}
		table := gmc.DB.Table("config")
		_, err = table.UpdateBy(gmap.M{"key": key}, gmap.M{"value": string(dStr)})
		if err != nil {
			this._JSONFail(err.Error())
		}
		this._JSONSuccess("")
	}
	value, err := this.configValue(key)
	if err != nil {
		this._JSONFail(err.Error())
	}
	this.View.SetMap(value)
	this.View.Layout("admin/form").Render("admin/config/" + key)
}
func (this *Config) ClearCache() {
	global.Context.Cache().Clear()
	this._JSONSuccess("清理成功")
}
