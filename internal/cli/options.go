package cli

type Options struct {
	ManifestPath string
	Preset       string
	Out          string
	Adapter      string
	DryRun       bool
	Force        bool
}
