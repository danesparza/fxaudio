package scripts

import "embed"

//go:embed sqlite/migrations/*
var FS embed.FS
