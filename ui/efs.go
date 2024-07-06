package ui

import "embed"

//go:embed "html" "static"
var Files embed.FS

//go:embed frontend/dist
var Frontend embed.FS
