package emd

import (
	"embed"
	_ "embed"
)

//go:embed assets/language.json
var ResLanguage string

//go:embed assets/web
var ResWeb embed.FS
