package global

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
)

var (
	Context *BContext
)

type BContext struct {
	db         gcore.Database
	cache      gcore.Cache
	log        gcore.Logger
	config     gcore.Config
	configFile string
	isDebug    bool
	server     gcore.HTTPServer
}

func (B *BContext) Server() gcore.HTTPServer {
	return B.server
}

func (B *BContext) SetServer(server gcore.HTTPServer) {
	B.server = server
}

func (B *BContext) ConfigFile() string {
	return B.configFile
}

func (B *BContext) Config() gcore.Config {
	return B.config
}

func (B *BContext) Log() gcore.Logger {
	return B.log
}

func (B *BContext) SetLog(log gcore.Logger) {
	B.log = log
}

func (B *BContext) Cache() gcore.Cache {
	return B.cache
}

func (B *BContext) SetCache(cache gcore.Cache) {
	B.cache = cache
}

func (B *BContext) DB() gcore.Database {
	return B.db
}

func (B *BContext) SetDB(db gcore.Database) {
	B.db = db
}
func (B *BContext) IsDebug() bool {
	return B.isDebug
}

func (B *BContext) SetIsDebug(isDebug bool) {
	B.isDebug = isDebug
}

func NewBContext(configFile string) (c *BContext, err error) {
	cfg, err := gconfig.NewConfigFile(configFile)
	if err != nil {
		return
	}
	c = &BContext{
		configFile: configFile,
		config:     cfg,
	}
	return
}