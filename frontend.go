package starter

import "embed"

//go:embed frontend/.output/public/**
var FrontendEmbeddedFiles embed.FS
