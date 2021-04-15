package assets

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS
