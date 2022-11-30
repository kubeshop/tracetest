package web

import "embed"

var (
	//go:embed build/*
	SPA embed.FS

	//go:embed build/static/*
	Static embed.FS
)
