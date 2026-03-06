package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultManifestDir      = "assets"
	defaultManifestFilename = "manifest.yml"
)

func ParseAndValidate() (Options, error) {
	opts, err := parse()
	if err != nil {
		return Options{}, err
	}

	if err := normalize(&opts); err != nil {
		return Options{}, err
	}

	if err := validate(opts); err != nil {
		return Options{}, err
	}

	return opts, nil
}

func parse() (Options, error) {
	defaultManifestPath := filepath.Join(defaultManifestDir, defaultManifestFilename)

	manifestHelp := fmt.Sprintf(
		"Path to the context manifest file (default: %s)",
		defaultManifestPath,
	)
	presetHelp := fmt.Sprintf(
		"Preset to use: %s (default: %s)",
		supportedPresetsText(),
		DefaultPreset,
	)

	manifestPath := flag.String("manifest", defaultManifestPath, manifestHelp)
	preset := flag.String("preset", DefaultPreset, presetHelp)
	out := flag.String("out", "", "output root directory where .context/ will be created (required)")
	dryRun := flag.Bool("dry-run", false, "preview actions without writing files (default: false)")
	force := flag.Bool("force", false, "overwrite existing destination files (default: false)")

	flag.Parse()

	return Options{
		ManifestPath: *manifestPath,
		Preset:       *preset,
		Out:          *out,
		DryRun:       *dryRun,
		Force:        *force,
	}, nil
}

func normalize(opts *Options) error {
	absManifestPath, err := filepath.Abs(opts.ManifestPath)
	if err != nil {
		return fmt.Errorf("failed to resolve manifest path: %w", err)
	}

	opts.ManifestPath = absManifestPath
	return nil
}

func validate(opts Options) error {
	if err := validateManifestPath(opts.ManifestPath); err != nil {
		return err
	}

	if err := validatePreset(opts.Preset); err != nil {
		return err
	}

	if opts.Out == "" {
		return fmt.Errorf("flag -out is required: specify the output root directory")
	}

	return nil
}

func validateManifestPath(path string) error {
	if path == "" {
		return fmt.Errorf("manifest path must not be empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("manifest file not found: %s", path)
		}
		return fmt.Errorf("failed to access manifest path %q: %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf(
			"manifest path points to a directory: %s (please provide a manifest file path, e.g. %s)",
			path,
			filepath.Join(path, defaultManifestFilename),
		)
	}

	return nil
}
