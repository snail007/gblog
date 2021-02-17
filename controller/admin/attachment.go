package admin

import (
	"fmt"
	gfile "github.com/snail007/gmc/util/file"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Attachment struct {
	Admin
}

func (this *Attachment) Upload() {
	if this.Ctx.IsPOST() {
		// do upload
		file, err := this.Ctx.FormFile("file", 0)
		if err != nil {
			this._JSONFail(err.Error())
		}
		rootDir := this.Config.GetString("attachment.dir")
		subDir := time.Now().Format("2006-01-02")
		uploadDir := filepath.Join(rootDir, subDir)
		if !gfile.Exists(uploadDir) {
			os.MkdirAll(uploadDir, 0755)
		}
		ext := filepath.Ext(file.Filename)
		if ext == "" {
			this._JSONFail("unknown file extension")
		}
		rid := fmt.Sprintf("%d_%d%s", time.Now().Unix(), rand.Int63(), ext)
		err = this.Ctx.SaveUploadedFile(file, filepath.Join(uploadDir, rid))
		if err != nil {
			this._JSONFail(err.Error())
		}
		urlPath := this.Config.GetString("attachment.url")
		this._JSONSuccess("", "", urlPath+"?id="+subDir+"/"+rid)
	} else {
		this.Ctx.WriteHeader(http.StatusBadRequest)
	}
}
