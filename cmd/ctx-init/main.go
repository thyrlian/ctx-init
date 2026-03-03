package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/thyrlian/ctx-init/internal/cli"
	"github.com/thyrlian/ctx-init/internal/manifest"
	"github.com/thyrlian/ctx-init/internal/plan"
	"github.com/thyrlian/ctx-init/internal/render"
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

	p, err := plan.Build(m, opts.ManifestPath, opts.Preset, ".", plan.Options{VerifySources: false})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("------------")
	fmt.Println("=== 🛠️ CTX-INIT 🛠️ ===")
	fmt.Printf("Manifest: %s\n", opts.ManifestPath)
	fmt.Printf("Preset:   %s\n", opts.Preset)
	fmt.Printf("version:      %d\n", m.Version)
	fmt.Printf("root_dir:     %s\n", m.RootDir)
	fmt.Printf("content_root: %s\n", m.ContentRoot)
	fmt.Printf("sections:     %d\n", len(m.Sections))
	fmt.Printf("presets:      %s\n", strings.Join(presetNames(m.Presets), ", "))
	fmt.Printf("files_total:  %d\n", countAllFiles(m.Sections))
	fmt.Printf("Plan:         %d entries (mode=%s)\n", len(p.Entries), p.Mode)
	for i, e := range p.Entries {
		fmt.Printf("%2d) %s\n", i+1, e.Src)
		fmt.Printf("        → %s\n", e.Dst)
	}

	fmt.Println()
	fmt.Println("Render:")
	result, err := render.Run(p, render.Options{
		DryRun: opts.DryRun,
		Force:  opts.Force,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	total := result.Copied + result.Skipped
	if opts.DryRun {
		fmt.Printf("Done (dry-run): %d total — %d would copy, %d would skip.\n", total, result.Copied, result.Skipped)
	} else {
		fmt.Printf("Done: %d total — %d copied, %d skipped.\n", total, result.Copied, result.Skipped)
	}
	fmt.Println("------------")
}

func presetNames(presets map[string]manifest.Preset) []string {
	names := make([]string, 0, len(presets))
	for name := range presets {
		names = append(names, name)
	}
	// Intentionally not sorted: keep original map iteration order for now
	return names
}

func countAllFiles(sections []manifest.Section) int {
	total := 0
	var walk func([]manifest.Section)
	walk = func(ss []manifest.Section) {
		for _, s := range ss {
			total += len(s.Files)
			if len(s.Children) > 0 {
				walk(s.Children)
			}
		}
	}
	walk(sections)
	return total
}
