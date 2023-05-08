package app

import (
	"embed"
)

//go:embed assets/*
var staticAssets embed.FS

func Assets() embed.FS {
	return staticAssets
}
