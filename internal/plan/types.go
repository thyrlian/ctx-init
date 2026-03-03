package plan

type Entry struct {
	// Absolute source path (template file)
	Src string

	// Absolute destination path (target project)
	Dst string

	// Effective tags = section tags (including parents) + file tags (deduped)
	Tags []string
}

type Plan struct {
	Preset  string
	Mode    string
	Entries []Entry
}

type Options struct {
	// If true, verify that each Src exists and is a file (double insurance).
	// Default should be false because manifest parser already validated sources.
	VerifySources bool
}
