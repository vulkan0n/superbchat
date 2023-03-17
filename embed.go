package embedfiles

import "embed"

//go:embed ui/html/*
var resources embed.FS

//go:embed ui/static
var styleFiles embed.FS

//go:embed config.json
var configBytes []byte
