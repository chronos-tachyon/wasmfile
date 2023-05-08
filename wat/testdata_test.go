package wat

import (
	"embed"
)

//go:embed testdata/*
var testDataFS embed.FS
