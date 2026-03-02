package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseFile(manifestPath string) (*Manifest, error) {
	manifestAbs, err := filepath.Abs(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("resolve manifest path %q: %w", manifestPath, err)
	}

	data, err := os.ReadFile(manifestAbs)
	if err != nil {
		return nil, fmt.Errorf("read manifest file %q: %w", manifestAbs, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest yaml %q: %w", manifestAbs, err)
	}

	manifestDir := filepath.Dir(manifestAbs)
	if err := validateBasic(&m, manifestDir); err != nil {
		return nil, err
	}

	return &m, nil
}

func validateBasic(m *Manifest, manifestDir string) error {
	if m.Version <= 0 {
		return fmt.Errorf("invalid manifest version: %d", m.Version)
	}
	if strings.TrimSpace(m.RootDir) == "" {
		return fmt.Errorf("root_dir must not be empty")
	}
	if strings.TrimSpace(m.ContentRoot) == "" {
		return fmt.Errorf("content_root must not be empty")
	}

	// content_root must be a relative path (portable templates)
	if filepath.IsAbs(m.ContentRoot) {
		return fmt.Errorf("content_root must be a relative path, got: %s", m.ContentRoot)
	}
	if strings.HasPrefix(m.ContentRoot, "~") {
		return fmt.Errorf("content_root must not use ~ expansion, got: %s", m.ContentRoot)
	}

	// Resolve content_root relative to the manifest directory for validation
	contentRootAbs := filepath.Clean(filepath.Join(manifestDir, m.ContentRoot))

	info, err := os.Stat(contentRootAbs)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("content_root not found: %s", contentRootAbs)
		}
		return fmt.Errorf("failed to access content_root %q: %w", contentRootAbs, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("content_root is not a directory: %s", contentRootAbs)
	}

	// Validate presets.mode (logic-bearing field)
	for name, p := range m.Presets {
		switch p.Mode {
		case ModeInclude:
			if len(p.Tags) == 0 {
				return fmt.Errorf("preset %q: mode=include requires non-empty tags", name)
			}
		case ModeFull:
			// ok, tags may be empty/omitted
		default:
			return fmt.Errorf("preset %q: invalid mode %q (supported: %s, %s)", name, p.Mode, ModeInclude, ModeFull)
		}
	}

	// Validate that all referenced source files exist under content_root
	for _, s := range m.Sections {
		if err := validateSectionSources(contentRootAbs, "", s); err != nil {
			return err
		}
	}

	return nil
}

func validateSectionSources(contentRootAbs, parentDir string, s Section) error {
	// Build the effective directory for this section
	secDir := strings.TrimSpace(s.Dir)
	effectiveDir := parentDir
	if secDir != "" {
		effectiveDir = filepath.Join(parentDir, secDir)
	}
	effectiveDir = filepath.Clean(effectiveDir)

	// Validate files in this section
	for _, f := range s.Files {
		name := strings.TrimSpace(f.Name)
		if name == "" {
			return fmt.Errorf("file.name must not be empty (section dir=%q)", s.Dir)
		}

		srcAbs := filepath.Clean(filepath.Join(contentRootAbs, effectiveDir, name))

		info, err := os.Stat(srcAbs)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("source file not found: %s", srcAbs)
			}
			return fmt.Errorf("failed to access source file %q: %w", srcAbs, err)
		}
		if info.IsDir() {
			return fmt.Errorf("source path is a directory, expected a file: %s", srcAbs)
		}
	}

	// Recurse children
	for _, c := range s.Children {
		if err := validateSectionSources(contentRootAbs, effectiveDir, c); err != nil {
			return err
		}
	}

	return nil
}
