package embedfiles

import "embed"

//go:embed ui/html/*
var Resources embed.FS

//go:embed ui/static/*
var StyleFiles embed.FS

//go:embed config.json
var ConfigBytes []byte