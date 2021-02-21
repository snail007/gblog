module gblog

go 1.12

require (
	github.com/google/go-github/v33 v33.0.0
	github.com/gookit/validate v1.2.8
	github.com/snail007/gmc v0.0.0-20210221072218-314f901975d6
	github.com/spf13/pflag v1.0.3
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

//replace github.com/snail007/gmc => ../gmc
