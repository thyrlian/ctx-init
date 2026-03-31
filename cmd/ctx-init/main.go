package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/thyrlian/ctx-init/internal/adapter"
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

	p, err := plan.Build(m, opts.ManifestPath, opts.Preset, opts.ProjectRoot, plan.Options{VerifySources: false})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("------------")
	fmt.Println("=== 🛠️ CTX-INIT 🛠️ ===")
	fmt.Printf("Manifest: %s\n", opts.ManifestPath)
	fmt.Printf("Preset:   %s\n", opts.Preset)
	if opts.Adapter != "" {
		fmt.Printf("Adapter:  %s\n", opts.Adapter)
	}
	fmt.Printf("version:      %d\n", m.Version)
	fmt.Printf("root_dir:     %s\n", m.RootDir)
	fmt.Printf("content_root: %s\n", m.ContentRoot)
	fmt.Printf("sections:     %d\n", len(m.Sections))
	fmt.Printf("presets:      %s\n", strings.Join(presetNames(m.Presets), ", "))
	fmt.Printf("files_total:  %d\n", countAllFiles(m.Sections))
	fmt.Printf("Plan:         %d entries (mode=%s)\n", len(p.Entries), p.Mode)
	for i, e := range p.Entries {
		fmt.Printf("%2d) %s\n", i+1, e.Src)
		fmt.Printf("        -> %s\n", e.Dst)
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

	var adapterResultPtr *adapter.Result
	if opts.Adapter != "" {
		fmt.Println()
		fmt.Println("Adapter:")
		adapterResult, err := adapter.Generate(opts.Adapter, opts.ProjectRoot, adapter.Options{
			DryRun: opts.DryRun,
			Force:  opts.Force,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		adapterResultPtr = &adapterResult
	}

	fmt.Println()
	printDoneSummary(result, adapterResultPtr, opts.DryRun)
	fmt.Println("------------")
}

func presetNames(presets map[string]manifest.Preset) []string {
	names := make([]string, 0, len(presets))
	for name := range presets {
		names = append(names, name)
	}
	sort.Strings(names)
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

func printDoneSummary(renderResult render.Result, adapterResult *adapter.Result, dryRun bool) {
	contextTotal := renderResult.Generated + renderResult.Copied + renderResult.Skipped

	if dryRun {
		fmt.Println("Done (dry-run):")
		fmt.Printf("  Context: %d total - %d would generate, %d would copy, %d would skip.\n", contextTotal, renderResult.Generated, renderResult.Copied, renderResult.Skipped)
		if adapterResult != nil {
			generated, skipped := adapterCounts(*adapterResult, true)
			fmt.Printf("  Adapter: %d would generate, %d would skip.\n", generated, skipped)
		}
		return
	}

	fmt.Println("Done:")
	fmt.Printf("  Context: %d total - %d generated, %d copied, %d skipped.\n", contextTotal, renderResult.Generated, renderResult.Copied, renderResult.Skipped)
	if adapterResult != nil {
		generated, skipped := adapterCounts(*adapterResult, false)
		fmt.Printf("  Adapter: %d generated, %d skipped.\n", generated, skipped)
	}
}

func adapterCounts(result adapter.Result, dryRun bool) (generated, skipped int) {
	if dryRun {
		switch result.Action {
		case adapter.ActionDryRunGenerate:
			return 1, 0
		case adapter.ActionDryRunSkip:
			return 0, 1
		default:
			return 0, 0
		}
	}

	switch result.Action {
	case adapter.ActionGenerated:
		return 1, 0
	case adapter.ActionSkipped:
		return 0, 1
	default:
		return 0, 0
	}
}
