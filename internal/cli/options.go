package cli

type Options struct {
	ManifestPath string
	Preset       string
	ProjectRoot  string
	Adapter      string
	DryRun       bool
	Force        bool
}
