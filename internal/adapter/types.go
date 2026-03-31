package adapter

import "io"

type Action string

const (
	ActionGenerated      Action = "generated"
	ActionSkipped        Action = "skipped"
	ActionDryRunGenerate Action = "dry-run/generate"
	ActionDryRunSkip     Action = "dry-run/skip"
)

type Options struct {
	DryRun bool
	Force  bool
	Writer io.Writer
}

type Result struct {
	Action        Action
	GeneratedPath string
	Skipped       bool
	UsedFallback  bool
	Message       string
}
