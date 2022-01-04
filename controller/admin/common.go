package admin

import (
	"encoding/base64"
	"strings"

	"gblog/global"
	"github.com/snail007/gmc"
	gcast "github.com/snail007/gmc/util/cast"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/xanzy/go-gitlab"
)

var (
	syncFolder = "articles"
)

func syncArticle(articleID string, oldArticle gmap.Mss) (err error) {
	if !isGitLabSync() {
		return
	}
	catalogTable := gmc.DB.Table("catalog")
	oldArticlePath := ""
	if len(oldArticle) > 0 {
		oldCatalog, err := catalogTable.GetByID(oldArticle["catalog_id"])
		if err != nil {
			global.Context.Log().Warnf("sync article error: %s", err)
			return err
		}
		oldArticlePath = syncFolder + "/" + oldCatalog["name"] + "/" + oldArticle["title"] + ".md"
	}
	articleTable := gmc.DB.Table("article")
	article, err := articleTable.GetByID(articleID)
	if err != nil {
		global.Context.Log().Warnf("sync article error: %s", err)
		return
	}
	catalog, err := catalogTable.GetByID(article["catalog_id"])
	if err != nil {
		global.Context.Log().Warnf("sync article error: %s", err)
		return
	}
	filePath := syncFolder + "/" + catalog["name"] + "/" + article["title"] + ".md"
	token := gcast.ToString(global.Context.BConfig("upload.gitlab_token"))
	apiURL := gcast.ToString(global.Context.BConfig("upload.gitlab_api_url"))
	userRepo := gcast.ToString(global.Context.BConfig("upload.gitlab_repo"))
	storageType := global.Context.BConfig("upload.upload_file_storage")
	if storageType != "gitlab" {
		return
	}
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL))
	if err != nil {
		global.Context.Log().Warnf("sync article error: %s", err)
		return
	}
	defaultBranchName, err := getDefaultBranchName()
	if err != nil {
		global.Context.Log().Warnf("sync article error: %s", err)
		return err
	}
	encoding := "base64"
	msg := "new file"
	binText := base64.StdEncoding.EncodeToString([]byte(article["content"]))
	_, _, err = client.RepositoryFiles.CreateFile(userRepo, filePath, &gitlab.CreateFileOptions{
		Branch:        &defaultBranchName,
		Encoding:      &encoding,
		Content:       &binText,
		CommitMessage: &msg,
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			// try update
			msg = "update file"
			_, _, err = client.RepositoryFiles.UpdateFile(userRepo, filePath, &gitlab.UpdateFileOptions{
				Branch:        &defaultBranchName,
				Encoding:      &encoding,
				Content:       &binText,
				CommitMessage: &msg,
			})
			if err != nil {
				global.Context.Log().Warnf("sync article error: %s", err)
				return
			}
			if oldArticlePath != "" && oldArticlePath != filePath {
				go func() {
					syncDeleteArticle(oldArticle)
				}()
			}
		} else {
			global.Context.Log().Warnf("sync article error: %s", err)
			return
		}
	} else {
		if oldArticlePath != "" && oldArticlePath != filePath {
			go func() {
				syncDeleteArticle(oldArticle)
			}()
		}
	}
	return
}

func syncDeleteArticle(article gmap.Mss) (err error) {
	isSync := gcast.ToBool(global.Context.BConfig("upload.gitlab_sync"))
	if !isSync {
		return
	}
	catalogTable := gmc.DB.Table("catalog")
	catalog, err := catalogTable.GetByID(article["catalog_id"])
	if err != nil {
		global.Context.Log().Warnf("sync delete article error: %s", err)
		return err
	}
	articlePath := syncFolder + "/" + catalog["name"] + "/" + article["title"] + ".md"
	client, err := getGitLabClient()
	if err != nil {
		global.Context.Log().Warnf("sync delete article error: %s", err)
		return err
	}

	defaultBranchName, err := getDefaultBranchName()
	if err != nil {
		global.Context.Log().Warnf("get default branch name error: %s", err)
		return err
	}
	msg := "delete file"
	resp, err := client.RepositoryFiles.DeleteFile(getUserRepo(), articlePath, &gitlab.DeleteFileOptions{
		Branch:        &defaultBranchName,
		CommitMessage: &msg,
	})
	resp.Body.Close()
	return
}

func getGitLabClient() (client *gitlab.Client, err error) {
	token := gcast.ToString(global.Context.BConfig("upload.gitlab_token"))
	apiURL := gcast.ToString(global.Context.BConfig("upload.gitlab_api_url"))
	storageType := global.Context.BConfig("upload.upload_file_storage")
	if storageType != "gitlab" {
		return
	}
	client, err = gitlab.NewClient(token, gitlab.WithBaseURL(apiURL))
	if err != nil {
		global.Context.Log().Warnf("new gitlab client error: %s", err)
		return
	}
	return
}

func getUserRepo() string {
	return gcast.ToString(global.Context.BConfig("upload.gitlab_repo"))
}

func isGitLabSync() bool {
	return gcast.ToBool(global.Context.BConfig("upload.gitlab_sync"))
}

func getDefaultBranchName() (string, error) {
	client, err := getGitLabClient()
	if err != nil {
		return "", err
	}
	bs, resp, err := client.Branches.ListBranches(getUserRepo(), &gitlab.ListBranchesOptions{})
	if err != nil {
		global.Context.Log().Warnf("get default branch name error: %s", err)
		return "", err
	}
	resp.Body.Close()
	defaultBranchName := ""
	for _, v := range bs {
		if v.Default {
			defaultBranchName = v.Name
			break
		}
	}
	return defaultBranchName, nil
}
