package adapter

import "io"

type Options struct {
	DryRun bool
	Force  bool
	Writer io.Writer
}

type Result struct {
	GeneratedPath string
	Skipped       bool
	UsedFallback  bool
	Message       string
}
