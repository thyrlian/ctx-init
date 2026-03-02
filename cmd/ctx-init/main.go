package main

import (
	"fmt"
	"os"

	"github.com/thyrlian/ctx-init/internal/cli"
	"github.com/thyrlian/ctx-init/internal/manifest"
)

func main() {
	opts, err := cli.ParseAndValidate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	m, err := manifest.ParseFile(opts.ManifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("------------")
	fmt.Println("=== 🛠️ CTX-INIT 🛠️ ===")
	fmt.Printf("Manifest: %s\n", opts.ManifestPath)
	fmt.Printf("Preset:   %s\n", opts.Preset)
	fmt.Println("------------")

	fmt.Println("------------")
	fmt.Printf("version:      %d\n", m.Version)
	fmt.Printf("root_dir:     %s\n", m.RootDir)
	fmt.Printf("content_root: %s\n", m.ContentRoot)
	fmt.Printf("sections:     %d\n", len(m.Sections))
	fmt.Printf("presets:      %d\n", len(m.Presets))
	fmt.Println("------------")
}
