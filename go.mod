module gblog

go 1.16

require (
	github.com/blevesearch/bleve v1.0.14
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/go-github/v33 v33.0.0
	github.com/gookit/validate v1.2.8
	github.com/snail007/gmc v0.0.0-20231031053155-291574b7f6a5
	github.com/snail007/resize v0.0.0-20180221191011-83c6a9932646
	github.com/snail007/uasurfer v0.0.0-20191028135549-26b5daa857f1
	github.com/spf13/pflag v1.0.3
	github.com/xanzy/go-gitlab v0.52.2
	github.com/yanyiwu/gojieba v1.1.2
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

//replace github.com/snail007/gmc => ../gmc
