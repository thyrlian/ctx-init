package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	presetMinimal  = "minimal"
	presetStandard = "standard"
	presetFull     = "full"

	defaultManifestDir      = "assets"
	defaultManifestFilename = "manifest.yml"
	defaultPreset           = presetStandard
)

var supportedPresets = []string{
	presetMinimal,
	presetStandard,
	presetFull,
}

type cliOptions struct {
	manifestPath string
	preset       string
}

func parseCLI() (cliOptions, error) {
	defaultManifestPath := filepath.Join(defaultManifestDir, defaultManifestFilename)

	manifestHelp := fmt.Sprintf(
		"Path to the context manifest file (default: %s)",
		defaultManifestPath,
	)
	presetHelp := fmt.Sprintf(
		"Preset to use: %s (default: %s)",
		strings.Join(supportedPresets, " | "),
		defaultPreset,
	)

	manifestPath := flag.String("manifest", defaultManifestPath, manifestHelp)
	preset := flag.String("preset", defaultPreset, presetHelp)

	flag.Parse()

	return cliOptions{
		manifestPath: *manifestPath,
		preset:       *preset,
	}, nil
}

func normalizeCLI(opts *cliOptions) error {
	absManifestPath, err := filepath.Abs(opts.manifestPath)
	if err != nil {
		return fmt.Errorf("failed to resolve manifest path: %w", err)
	}

	opts.manifestPath = absManifestPath
	return nil
}

func validateCLI(opts cliOptions) error {
	if err := validateManifestPath(opts.manifestPath); err != nil {
		return err
	}

	if err := validatePreset(opts.preset); err != nil {
		return err
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

func validatePreset(preset string) error {
	for _, p := range supportedPresets {
		if preset == p {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid preset %q (supported: %s)",
		preset,
		strings.Join(supportedPresets, ", "),
	)
}

func main() {
	opts, err := parseCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to parse CLI args: %v\n", err)
		os.Exit(1)
	}

	if err := normalizeCLI(&opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := validateCLI(opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("------------")
	fmt.Println("ctx-init (skeleton)")
	fmt.Println("------------")

	fmt.Printf("Manifest: %s\n", opts.manifestPath)
	fmt.Printf("Preset:   %s\n", opts.preset)
}
