module gblog

go 1.12

require (
	github.com/google/go-github/v33 v33.0.0
	github.com/gookit/validate v1.2.8
	github.com/snail007/gmc v0.0.0-20210222085810-f2628a04c9da
	github.com/snail007/resize v0.0.0-20180221191011-83c6a9932646
	github.com/spf13/pflag v1.0.3
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

//replace github.com/snail007/gmc => ../gmc
