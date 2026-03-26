package plan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/thyrlian/ctx-init/internal/manifest"
)

// ---- mergeTags -------------------------------------------------------

func TestMergeTags(t *testing.T) {
	tests := []struct {
		name string
		a, b []string
		want []string
	}{
		{
			name: "both empty",
			a:    nil, b: nil,
			want: []string{},
		},
		{
			name: "a only",
			a:    []string{"core", "global"},
			want: []string{"core", "global"},
		},
		{
			name: "b only",
			b:    []string{"vcs", "git"},
			want: []string{"vcs", "git"},
		},
		{
			name: "no overlap",
			a:    []string{"core"},
			b:    []string{"vcs"},
			want: []string{"core", "vcs"},
		},
		{
			name: "with overlap — deduped",
			a:    []string{"core", "global"},
			b:    []string{"global", "vcs"},
			want: []string{"core", "global", "vcs"},
		},
		{
			name: "whitespace trimmed and deduped",
			a:    []string{" core ", "global"},
			b:    []string{"core", " global "},
			want: []string{"core", "global"},
		},
		{
			name: "empty strings stripped",
			a:    []string{"core", ""},
			b:    []string{"", "vcs"},
			want: []string{"core", "vcs"},
		},
		{
			name: "a-order preserved",
			a:    []string{"z", "a", "m"},
			b:    []string{"b"},
			want: []string{"z", "a", "m", "b"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeTags(tc.a, tc.b)
			// nil result should equal empty expected
			if len(got) == 0 && len(tc.want) == 0 {
				return
			}
			if len(got) != len(tc.want) {
				t.Fatalf("mergeTags(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("mergeTags(%v, %v)[%d] = %q; want %q", tc.a, tc.b, i, got[i], tc.want[i])
				}
			}
		})
	}
}

// ---- intersects -------------------------------------------------------

func TestIntersects(t *testing.T) {
	selected := map[string]struct{}{
		"core":   {},
		"global": {},
	}

	tests := []struct {
		name string
		tags []string
		want bool
	}{
		{"empty tags", nil, false},
		{"no overlap", []string{"vcs", "git"}, false},
		{"single match", []string{"core"}, true},
		{"match at end", []string{"vcs", "global"}, true},
		{"all match", []string{"core", "global"}, true},
		{"case sensitive — no match", []string{"Core"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := intersects(tc.tags, selected); got != tc.want {
				t.Errorf("intersects(%v) = %v; want %v", tc.tags, got, tc.want)
			}
		})
	}
}

// ---- filterByPreset ---------------------------------------------------

func TestFilterByPreset(t *testing.T) {
	entries := []Entry{
		{Dst: "a", Tags: []string{"core", "global"}},
		{Dst: "b", Tags: []string{"workflow"}},
		{Dst: "c", Tags: []string{"core"}},
		{Dst: "d", Tags: []string{"standards"}},
		{Dst: "e", Tags: []string{}},
	}

	t.Run("full mode returns all entries", func(t *testing.T) {
		preset := manifest.Preset{Mode: manifest.ModeFull}
		got := filterByPreset(entries, preset)
		if len(got) != len(entries) {
			t.Errorf("full mode: got %d entries; want %d", len(got), len(entries))
		}
	})

	t.Run("include mode — matches by tag", func(t *testing.T) {
		preset := manifest.Preset{Mode: manifest.ModeInclude, Tags: []string{"core"}}
		got := filterByPreset(entries, preset)
		if len(got) != 2 {
			t.Fatalf("include core: got %d entries; want 2", len(got))
		}
		for _, e := range got {
			if !intersects(e.Tags, map[string]struct{}{"core": {}}) {
				t.Errorf("entry %q should not have been included", e.Dst)
			}
		}
	})

	t.Run("include mode — no matching tags returns empty", func(t *testing.T) {
		preset := manifest.Preset{Mode: manifest.ModeInclude, Tags: []string{"nonexistent"}}
		got := filterByPreset(entries, preset)
		if len(got) != 0 {
			t.Errorf("expected 0 entries; got %d", len(got))
		}
	})

	t.Run("include mode — multiple tags act as OR", func(t *testing.T) {
		preset := manifest.Preset{Mode: manifest.ModeInclude, Tags: []string{"workflow", "standards"}}
		got := filterByPreset(entries, preset)
		if len(got) != 2 {
			t.Errorf("expected 2 entries; got %d", len(got))
		}
	})

	t.Run("include mode — empty preset tags returns empty", func(t *testing.T) {
		preset := manifest.Preset{Mode: manifest.ModeInclude, Tags: nil}
		got := filterByPreset(entries, preset)
		if len(got) != 0 {
			t.Errorf("expected 0 entries; got %d", len(got))
		}
	})

	t.Run("unknown mode returns nil", func(t *testing.T) {
		preset := manifest.Preset{Mode: "bogus"}
		got := filterByPreset(entries, preset)
		if got != nil {
			t.Errorf("expected nil; got %v", got)
		}
	})
}

// ---- Build ------------------------------------------------------------

func TestBuild_nilManifest(t *testing.T) {
	_, err := Build(nil, "manifest.yml", "standard", "/tmp/out", Options{})
	if err == nil {
		t.Fatal("expected error for nil manifest; got nil")
	}
}

func TestBuild_emptyPreset(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "assets/context",
		Presets:     map[string]manifest.Preset{"standard": {Mode: manifest.ModeFull}},
	}
	_, err := Build(m, "manifest.yml", "", "/tmp/out", Options{})
	if err == nil {
		t.Fatal("expected error for empty preset; got nil")
	}
}

func TestBuild_unknownPreset(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "assets/context",
		Presets:     map[string]manifest.Preset{"standard": {Mode: manifest.ModeFull}},
	}
	_, err := Build(m, "manifest.yml", "nonexistent", "/tmp/out", Options{})
	if err == nil {
		t.Fatal("expected error for unknown preset; got nil")
	}
}

func TestBuild_emptyManifestPath(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "assets/context",
		Presets:     map[string]manifest.Preset{"standard": {Mode: manifest.ModeFull}},
	}
	_, err := Build(m, "  ", "standard", "/tmp/out", Options{})
	if err == nil {
		t.Fatal("expected error for empty manifestPath; got nil")
	}
}

func TestBuild_emptyTargetDir(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "assets/context",
		Presets:     map[string]manifest.Preset{"standard": {Mode: manifest.ModeFull}},
	}
	_, err := Build(m, "manifest.yml", "standard", "  ", Options{})
	if err == nil {
		t.Fatal("expected error for empty targetDir; got nil")
	}
}

func TestBuild_fullPresetIncludesAllEntries(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "testdata",
		Sections: []manifest.Section{
			{
				Tags: []string{"core"},
				Files: []manifest.File{
					{Name: "ai_protocol.md", Tags: []string{}},
				},
			},
			{
				Dir:  "standards",
				Tags: []string{"standards"},
				Files: []manifest.File{
					{Name: "naming.md"},
					{Name: "security.md"},
				},
			},
		},
		Presets: map[string]manifest.Preset{
			"full": {Mode: manifest.ModeFull},
		},
	}

	p, err := Build(m, "manifest.yml", "full", "/tmp/out", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Entries) != 3 {
		t.Fatalf("full preset: got %d entries; want 3", len(p.Entries))
	}
	if p.Preset != "full" {
		t.Errorf("Preset = %q; want %q", p.Preset, "full")
	}
	if p.Mode != manifest.ModeFull {
		t.Errorf("Mode = %q; want %q", p.Mode, manifest.ModeFull)
	}
	// entries are sorted by Dst; verify Src and Dst path suffixes in order
	wantDstSuffixes := []string{
		filepath.Join(".context", "ai_protocol.md"),
		filepath.Join(".context", "standards", "naming.md"),
		filepath.Join(".context", "standards", "security.md"),
	}
	wantSrcSuffixes := []string{
		filepath.Join("testdata", "ai_protocol.md"),
		filepath.Join("testdata", "standards", "naming.md"),
		filepath.Join("testdata", "standards", "security.md"),
	}
	for i := range wantDstSuffixes {
		if !strings.HasSuffix(p.Entries[i].Dst, wantDstSuffixes[i]) {
			t.Errorf("Entries[%d].Dst = %q; want suffix %q", i, p.Entries[i].Dst, wantDstSuffixes[i])
		}
		if !strings.HasSuffix(p.Entries[i].Src, wantSrcSuffixes[i]) {
			t.Errorf("Entries[%d].Src = %q; want suffix %q", i, p.Entries[i].Src, wantSrcSuffixes[i])
		}
	}
}

func TestBuild_includePresetFiltersEntries(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "testdata",
		Sections: []manifest.Section{
			{
				Tags: []string{"core"},
				Files: []manifest.File{
					{Name: "ai_protocol.md"},
				},
			},
			{
				Dir:  "standards",
				Tags: []string{"standards"},
				Files: []manifest.File{
					{Name: "naming.md"},
					{Name: "security.md"},
				},
			},
		},
		Presets: map[string]manifest.Preset{
			"minimal": {Mode: manifest.ModeInclude, Tags: []string{"core"}},
		},
	}

	p, err := Build(m, "manifest.yml", "minimal", "/tmp/out", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Entries) != 1 {
		t.Errorf("minimal preset: got %d entries; want 1", len(p.Entries))
	}
	if len(p.Entries) > 0 {
		got := p.Entries[0].Tags
		if !intersects(got, map[string]struct{}{"core": {}}) {
			t.Errorf("entry tags %v do not include 'core'", got)
		}
	}
}

func TestBuild_sectionTagsInheritedByFiles(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "testdata",
		Sections: []manifest.Section{
			{
				Dir:  "standards",
				Tags: []string{"standards", "global"},
				Files: []manifest.File{
					{Name: "naming.md", Tags: []string{"extra"}},
				},
			},
		},
		Presets: map[string]manifest.Preset{
			"full": {Mode: manifest.ModeFull},
		},
	}

	p, err := Build(m, "manifest.yml", "full", "/tmp/out", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Entries) != 1 {
		t.Fatalf("expected 1 entry; got %d", len(p.Entries))
	}
	tags := p.Entries[0].Tags
	expected := []string{"standards", "global", "extra"}
	if len(tags) != len(expected) {
		t.Fatalf("tags = %v; want %v", tags, expected)
	}
	for i, want := range expected {
		if tags[i] != want {
			t.Errorf("tags[%d] = %q; want %q", i, tags[i], want)
		}
	}
}

func TestBuild_entriesSortedDeterministically(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "testdata",
		Sections: []manifest.Section{
			{
				Files: []manifest.File{
					{Name: "z_file.md"},
					{Name: "a_file.md"},
					{Name: "m_file.md"},
				},
			},
		},
		Presets: map[string]manifest.Preset{
			"full": {Mode: manifest.ModeFull},
		},
	}

	p, err := Build(m, "manifest.yml", "full", "/tmp/out", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 1; i < len(p.Entries); i++ {
		if p.Entries[i].Dst < p.Entries[i-1].Dst {
			t.Errorf("entries not sorted: %q before %q", p.Entries[i-1].Dst, p.Entries[i].Dst)
		}
	}
}

func TestBuild_skipsEmptyFileNames(t *testing.T) {
	m := &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: "testdata",
		Sections: []manifest.Section{
			{
				Files: []manifest.File{
					{Name: "valid.md"},
					{Name: ""},
					{Name: "  "},
				},
			},
		},
		Presets: map[string]manifest.Preset{
			"full": {Mode: manifest.ModeFull},
		},
	}

	p, err := Build(m, "manifest.yml", "full", "/tmp/out", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Entries) != 1 {
		t.Errorf("expected 1 entry (empty names skipped); got %d", len(p.Entries))
	}
}

// ---- VerifySources ----------------------------------------------------

func makeVerifySourcesManifest(contentRoot string) *manifest.Manifest {
	return &manifest.Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: contentRoot,
		Sections: []manifest.Section{
			{
				Dir: "standards",
				Files: []manifest.File{
					{Name: "naming.md"},
				},
			},
		},
		Presets: map[string]manifest.Preset{
			"full": {Mode: manifest.ModeFull},
		},
	}
}

func TestBuild_verifySources_fileExists(t *testing.T) {
	tmp := t.TempDir()

	srcDir := filepath.Join(tmp, "content", "standards")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "naming.md"), []byte("# Naming"), 0o644); err != nil {
		t.Fatal(err)
	}

	manifestPath := filepath.Join(tmp, "manifest.yml")
	m := makeVerifySourcesManifest("content")

	_, err := Build(m, manifestPath, "full", filepath.Join(tmp, "out"), Options{VerifySources: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_verifySources_fileMissing(t *testing.T) {
	tmp := t.TempDir()
	manifestPath := filepath.Join(tmp, "manifest.yml")
	m := makeVerifySourcesManifest("content") // content dir never created

	_, err := Build(m, manifestPath, "full", filepath.Join(tmp, "out"), Options{VerifySources: true})
	if err == nil {
		t.Fatal("expected error for missing source file; got nil")
	}
	if !strings.Contains(err.Error(), "source file check failed") {
		t.Errorf("error message %q should contain %q", err.Error(), "source file check failed")
	}
}

func TestBuild_verifySources_pathIsDir(t *testing.T) {
	tmp := t.TempDir()

	// create a directory where a file is expected
	dirPath := filepath.Join(tmp, "content", "standards", "naming.md")
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		t.Fatal(err)
	}

	manifestPath := filepath.Join(tmp, "manifest.yml")
	m := makeVerifySourcesManifest("content")

	_, err := Build(m, manifestPath, "full", filepath.Join(tmp, "out"), Options{VerifySources: true})
	if err == nil {
		t.Fatal("expected error when source path is a directory; got nil")
	}
	if !strings.Contains(err.Error(), "expected file, got directory") {
		t.Errorf("error = %q; want message containing %q", err.Error(), "expected file, got directory")
	}
}
