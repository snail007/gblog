// +build bleve

package initialize

import (
	"fmt"
	_ "gblog/util/bleve"
	"gblog/global"
	"gblog/util/bleve/dict"
	"github.com/blevesearch/bleve"
	gdb "github.com/snail007/gmc/module/db"
	"os"
	"path/filepath"
)

func initIndexer(ctx *global.BContext) (err error) {
	dictPath := "data/dict"
	os.RemoveAll(dictPath)
	err = os.Mkdir(dictPath, 0755)
	if err != nil {
		return
	}
	files, err := dict.Dict.ReadDir(".")
	if err != nil {
		return
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		pathD := filepath.Join(dictPath, filepath.Base(f.Name()))
		bs, err := dict.Dict.ReadFile(f.Name())
		if err != nil {
			return err
		}
		err = os.WriteFile(pathD, bs, 0644)
		if err != nil {
			return err
		}
	}
	indexMapping := bleve.NewIndexMapping()
	err = indexMapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dictpath":      filepath.Join(dictPath, "jieba.dict.utf8"),
			"hmmpath":       filepath.Join(dictPath, "hmm_model.utf8"),
			"userdictpath":  filepath.Join(dictPath, "user.dict.utf8"),
			"idfpath":       filepath.Join(dictPath, "idf.utf8"),
			"stopwordspath": filepath.Join(dictPath, "stop_words.utf8"),
			"type":          "gojieba",
		},
	)
	if err != nil {
		return
	}
	err = indexMapping.AddCustomAnalyzer("gojieba",
		map[string]interface{}{
			"type":      "gojieba",
			"tokenizer": "gojieba",
		},
	)
	if err != nil {
		return
	}
	indexMapping.DefaultAnalyzer = "gojieba"

	indexer, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return
	}
	ctx.SetIndexer(indexer)
	article := gdb.Table("article")
	rows, err := article.GetAll()
	ctx.Log().Infof("indexing articles ...")
	for _, row := range rows {
		doc := fmt.Sprintf("%s\n%s\n%s", row["title"], row["summary"], row["content"])
		err = indexer.Index(row["article_id"], doc)
		if err != nil {
			return
		}
	}
	ctx.Log().Infof("indexing articles success")
	return
}
