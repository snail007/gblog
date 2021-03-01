package bleve

import (
	"errors"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"github.com/yanyiwu/gojieba"
)

type JiebaTokenizer struct {
	handle *gojieba.Jieba
}

func newJiebaTokenizer(dictpath, hmmpath, userdictpath, idfpath, stopwordspath string) *JiebaTokenizer {
	x := gojieba.NewJieba(dictpath, hmmpath, userdictpath, idfpath, stopwordspath)
	return &JiebaTokenizer{x}
}

func (x *JiebaTokenizer) Free() {
	x.handle.Free()
}

func (x *JiebaTokenizer) Tokenize(sentence []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	pos := 1
	words := x.handle.Tokenize(string(sentence), gojieba.SearchMode, true)
	for _, word := range words {
		token := analysis.Token{
			Term:     []byte(word.Str),
			Start:    word.Start,
			End:      word.End,
			Position: pos,
			Type:     analysis.Ideographic,
		}
		result = append(result, &token)
		pos++
	}
	return result
}

func tokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	dictpath, ok := config["dictpath"].(string)
	if !ok {
		return nil, errors.New("config dictpath not found")
	}
	hmmpath, ok := config["hmmpath"].(string)
	if !ok {
		return nil, errors.New("config hmmpath not found")
	}
	userdictpath, ok := config["userdictpath"].(string)
	if !ok {
		return nil, errors.New("config userdictpath not found")
	}
	idfpath, ok := config["idfpath"].(string)
	if !ok {
		return nil, errors.New("config idfpath not found")
	}
	stopwordspath, ok := config["idfpath"].(string)
	if !ok {
		return nil, errors.New("config stopwordspath not found")
	}
	return newJiebaTokenizer(dictpath, hmmpath, userdictpath, idfpath, stopwordspath), nil
}

func init() {
	registry.RegisterTokenizer("gojieba", tokenizerConstructor)
}
