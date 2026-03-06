package plan

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/thyrlian/ctx-init/internal/manifest"
)

func Build(m *manifest.Manifest, manifestPath string, presetName string, targetDir string, opt Options) (*Plan, error) {
	if m == nil {
		return nil, fmt.Errorf("manifest must not be nil")
	}

	presetName = strings.TrimSpace(presetName)
	if presetName == "" {
		return nil, fmt.Errorf("preset name must not be empty")
	}

	if strings.TrimSpace(manifestPath) == "" {
		return nil, fmt.Errorf("manifestPath must not be empty")
	}
	if strings.TrimSpace(targetDir) == "" {
		return nil, fmt.Errorf("targetDir must not be empty")
	}

	preset, ok := m.Presets[presetName]
	if !ok {
		return nil, fmt.Errorf("preset %q not found in manifest", presetName)
	}

	manifestAbs, err := filepath.Abs(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("resolve manifest path %q: %w", manifestPath, err)
	}
	manifestDir := filepath.Dir(manifestAbs)

	// content_root is relative (validated by manifest parser)
	contentRootAbs := filepath.Clean(filepath.Join(manifestDir, m.ContentRoot))

	targetAbs, err := filepath.Abs(targetDir)
	if err != nil {
		return nil, fmt.Errorf("resolve targetDir %q: %w", targetDir, err)
	}
	outRootAbs := filepath.Clean(filepath.Join(targetAbs, m.RootDir))

	var entries []Entry
	for _, s := range m.Sections {
		entries = append(entries, buildEntriesForSection(contentRootAbs, outRootAbs, "", nil, s)...)
	}

	entries = filterByPreset(entries, preset)

	// Deterministic ordering (stable for tests and logs)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Dst == entries[j].Dst {
			return entries[i].Src < entries[j].Src
		}
		return entries[i].Dst < entries[j].Dst
	})

	if opt.VerifySources {
		for _, e := range entries {
			if err := requireFileExists(e.Src); err != nil {
				return nil, fmt.Errorf("source file check failed: %w", err)
			}
		}
	}

	return &Plan{
		Preset:     presetName,
		Mode:       preset.Mode,
		Entries:    entries,
		OutRootAbs: outRootAbs,
	}, nil
}

func buildEntriesForSection(contentRootAbs, outRootAbs, parentDir string, parentTags []string, s manifest.Section) []Entry {
	secDir := strings.TrimSpace(s.Dir)

	// logical dir: always forward-slash, OS-independent
	effectiveDir := parentDir
	if secDir != "" {
		effectiveDir = path.Join(parentDir, secDir)
	}
	effectiveDir = path.Clean(effectiveDir)
	if effectiveDir == "." {
		effectiveDir = ""
	}

	secTags := mergeTags(parentTags, s.Tags)

	var entries []Entry

	for _, f := range s.Files {
		name := strings.TrimSpace(f.Name)
		if name == "" {
			continue
		}

		src := filepath.Clean(filepath.Join(contentRootAbs, filepath.FromSlash(effectiveDir), name))
		dst := filepath.Clean(filepath.Join(outRootAbs, filepath.FromSlash(effectiveDir), name))

		entries = append(entries, Entry{
			Src:  src,
			Dst:  dst,
			Tags: mergeTags(secTags, f.Tags),
		})
	}

	for _, c := range s.Children {
		entries = append(entries, buildEntriesForSection(contentRootAbs, outRootAbs, effectiveDir, secTags, c)...)
	}

	return entries
}

func filterByPreset(entries []Entry, preset manifest.Preset) []Entry {
	switch preset.Mode {
	case manifest.ModeFull:
		return entries

	case manifest.ModeInclude:
		want := make(map[string]struct{}, len(preset.Tags))
		for _, t := range preset.Tags {
			t = strings.TrimSpace(t)
			if t != "" {
				want[t] = struct{}{}
			}
		}

		selectedEntries := make([]Entry, 0, len(entries))
		for _, e := range entries {
			if intersects(e.Tags, want) {
				selectedEntries = append(selectedEntries, e)
			}
		}
		return selectedEntries

	default:
		// parser validates; safe fallback
		return nil
	}
}

// intersects returns true if any tag in tags is present in want
func intersects(tags []string, want map[string]struct{}) bool {
	for _, t := range tags {
		if _, ok := want[t]; ok {
			return true
		}
	}
	return false
}

// mergeTags merges two tag slices into one, trimming spaces and deduping
func mergeTags(a []string, b []string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	out := make([]string, 0, len(a)+len(b))

	add := func(x string) {
		x = strings.TrimSpace(x)
		if x == "" {
			return
		}
		if _, ok := seen[x]; ok {
			return
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}

	for _, t := range a {
		add(t)
	}
	for _, t := range b {
		add(t)
	}
	return out
}

func requireFileExists(p string) error {
	info, err := os.Stat(p)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("expected file, got directory")
	}
	return nil
}
