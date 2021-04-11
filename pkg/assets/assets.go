package assets

import (
	"embed"
)

//go:embed tmpl/*
var Templates embed.FS
