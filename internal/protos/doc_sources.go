package protos

import "embed"

// DocumentationSources keeps protocol comments available in standalone binaries.
//
//go:embed *.go
var DocumentationSources embed.FS
