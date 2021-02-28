package dict

import "embed"

var (
	//go:embed *.utf8
	Dict embed.FS
)
