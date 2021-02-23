package admin

import (
	"context"
	"fmt"
	"gblog/global"
	"github.com/google/go-github/v33/github"
	gcast "github.com/snail007/gmc/util/cast"
	gfile "github.com/snail007/gmc/util/file"
	gmap "github.com/snail007/gmc/util/map"
	"golang.org/x/oauth2"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Attachment struct {
	Admin
}

func (this *Attachment) Upload() {
	if this.Ctx.IsPOST() {
		// do upload
		file, err := this.Ctx.FormFile("file", 0)
		isEditor := false
		if err != nil {
			isEditor = true
			file, err = this.Ctx.FormFile("editormd-image-file", 0)
		}
		if err != nil {
			this.jsonFail(isEditor, err.Error())
		}
		ext := filepath.Ext(file.Filename)
		if ext == "" {
			this.jsonFail(isEditor, "unknown file extension")
		}
		subDir := time.Now().Format("2006-01-02")
		randFilename := fmt.Sprintf("%d_%d%s", time.Now().Unix(), rand.Int63(), ext)

		storageType := global.Context.BConfig("upload.upload_file_storage")
		switch storageType {
		case "local":
			rootDir := this.Config.GetString("attachment.dir")
			uploadDir := filepath.Join(rootDir, subDir)
			if !gfile.Exists(uploadDir) {
				os.MkdirAll(uploadDir, 0755)
			}
			savePath := filepath.Join(uploadDir, randFilename)
			err = this.uploadToLocal(savePath, file)
		case "github":
			savePath := "attachment/" + subDir + "/" + randFilename
			err = this.uploadToGithub(savePath, file)
		}

		if err != nil {
			this.jsonFail(isEditor, err.Error())
		}

		urlPath := this.Config.GetString("attachment.url")
		link := urlPath + "/" + subDir + "/" + randFilename
		this.jsonSuccess(isEditor, "", link)
	} else {
		this.Ctx.WriteHeader(http.StatusBadRequest)
	}
}

func (this *Attachment) uploadToLocal(savePath string, file *multipart.FileHeader) (err error) {
	err = this.Ctx.SaveUploadedFile(file, savePath)
	return
}

func (this *Attachment) uploadToGithub(filePath string, file *multipart.FileHeader) (err error) {
	token := gcast.ToString(global.Context.BConfig("upload.github_token"))
	userRepo := gcast.ToString(global.Context.BConfig("upload.github_repo"))
	f, err := file.Open()
	if err != nil {
		return
	}
	defer f.Close()
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	data := strings.Split(userRepo, "/")
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel1()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx1, ts)
	client := github.NewClient(tc)
	repo, _, err := client.Repositories.Get(ctx1, data[0], data[1])
	if err != nil {
		return
	}
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("upload"),
		Content:   contents,
		Branch:    repo.DefaultBranch,
		Committer: &github.CommitAuthor{Name: github.String("gblog"), Email: github.String("bot@gblog")},
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel2()
	_, _, err = client.Repositories.CreateFile(ctx2, data[0], data[1], filePath, opts)
	return
}

func (this *Attachment) jsonSuccess(isEditor bool, msg string, data ...interface{}) {
	if isEditor {
		this.Ctx.JSON(200, gmap.M{
			"success": 1,
			"message": msg,
			"url":     data[0],
		})
	} else {
		this._JSONSuccess(msg, "", data[0])
	}
}
func (this *Attachment) jsonFail(isEditor bool, msg string) {
	if isEditor {
		this.Ctx.JSON(200, gmap.M{
			"success": 0,
			"message": msg,
			"url":     "",
		})
	} else {
		this._JSONSuccess(msg)
	}
}
