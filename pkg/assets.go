package pkg

import "embed"

//go:embed assets/*
var AssetData embed.FS
