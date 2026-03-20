package manifest

type Manifest struct {
	Version     int               `yaml:"version"`
	RootDir     string            `yaml:"root_dir"`
	ContentRoot string            `yaml:"content_root"`
	Sections    []Section         `yaml:"sections"`
	Presets     map[string]Preset `yaml:"presets"`
}

type Section struct {
	Dir      string    `yaml:"dir"`
	Tags     []string  `yaml:"tags"`
	Files    []File    `yaml:"files"`
	Children []Section `yaml:"children"`
}

type File struct {
	Name    string   `yaml:"name"`
	Tags    []string `yaml:"tags"`
	Gateway bool     `yaml:"gateway"`
}

type Preset struct {
	Mode string   `yaml:"mode"`
	Tags []string `yaml:"tags"`
}

const (
	ModeInclude = "include"
	ModeFull    = "full"
)
