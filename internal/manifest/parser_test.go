package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

// makeContentRoot creates a temp dir with a content subdirectory and returns both.
func makeContentRoot(t *testing.T) (manifestDir, contentDir string) {
	t.Helper()
	tmp := t.TempDir()
	content := filepath.Join(tmp, "content")
	if err := os.MkdirAll(content, 0o755); err != nil {
		t.Fatal(err)
	}
	return tmp, content
}

// minimalManifest returns a valid Manifest pointing at the given contentRoot name
// (relative to manifestDir).
func minimalManifest(contentRoot string) *Manifest {
	return &Manifest{
		Version:     1,
		RootDir:     ".context",
		ContentRoot: contentRoot,
		Presets:     map[string]Preset{"full": {Mode: ModeFull}},
	}
}

// ---- validateBasic ----------------------------------------------------

func TestValidateBasic_valid(t *testing.T) {
	manifestDir, _ := makeContentRoot(t)
	m := minimalManifest("content")
	if err := validateBasic(m, manifestDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateBasic_invalidVersion(t *testing.T) {
	manifestDir, _ := makeContentRoot(t)
	m := minimalManifest("content")
	m.Version = 0
	if err := validateBasic(m, manifestDir); err == nil {
		t.Fatal("expected error for version 0; got nil")
	}
}

func TestValidateBasic_emptyRootDir(t *testing.T) {
	manifestDir, _ := makeContentRoot(t)
	m := minimalManifest("content")
	m.RootDir = "  "
	if err := validateBasic(m, manifestDir); err == nil {
		t.Fatal("expected error for empty root_dir; got nil")
	}
}

func TestValidateBasic_emptyContentRoot(t *testing.T) {
	manifestDir, _ := makeContentRoot(t)
	m := minimalManifest("content")
	m.ContentRoot = "  "
	if err := validateBasic(m, manifestDir); err == nil {
		t.Fatal("expected error for empty content_root; got nil")
	}
}

func TestValidateBasic_absoluteContentRoot(t *testing.T) {
	manifestDir, _ := makeContentRoot(t)
	m := minimalManifest("/absolute/path")
	if err := validateBasic(m, manifestDir); err == nil {
		t.Fatal("expected error for absolute content_root; got nil")
	}
}

func TestValidateBasic_tildeContentRoot(t *testing.T) {
	manifestDir, _ := makeContentRoot(t)
	m := minimalManifest("~/content")
	if err := validateBasic(m, manifestDir); err == nil {
		t.Fatal("expected error for ~ content_root; got nil")
	}
}

func TestValidateBasic_contentRootNotFound(t *testing.T) {
	tmp := t.TempDir() // no "content" subdir created
	m := minimalManifest("content")
	if err := validateBasic(m, tmp); err == nil {
		t.Fatal("expected error for missing content_root; got nil")
	}
}

func TestValidateBasic_contentRootIsFile(t *testing.T) {
	tmp := t.TempDir()
	// create a file where a directory is expected
	filePath := filepath.Join(tmp, "content")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := minimalManifest("content")
	if err := validateBasic(m, tmp); err == nil {
		t.Fatal("expected error when content_root is a file; got nil")
	}
}

func TestValidateBasic_includePresetRequiresTags(t *testing.T) {
	manifestDir, _ := makeContentRoot(t)
	m := minimalManifest("content")
	m.Presets["minimal"] = Preset{Mode: ModeInclude, Tags: nil}
	if err := validateBasic(m, manifestDir); err == nil {
		t.Fatal("expected error for include preset with no tags; got nil")
	}
}

func TestValidateBasic_unknownPresetMode(t *testing.T) {
	manifestDir, _ := makeContentRoot(t)
	m := minimalManifest("content")
	m.Presets["bogus"] = Preset{Mode: "all"}
	if err := validateBasic(m, manifestDir); err == nil {
		t.Fatal("expected error for unknown preset mode; got nil")
	}
}

// ---- validateSectionSources -------------------------------------------

func TestValidateSectionSources_valid(t *testing.T) {
	_, contentDir := makeContentRoot(t)
	if err := os.MkdirAll(filepath.Join(contentDir, "standards"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "standards", "naming.md"), []byte("# Naming"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := Section{
		Dir:   "standards",
		Files: []File{{Name: "naming.md"}},
	}
	if err := validateSectionSources(contentDir, "", s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSectionSources_nestedChildValid(t *testing.T) {
	_, contentDir := makeContentRoot(t)
	if err := os.MkdirAll(filepath.Join(contentDir, "standards", "go"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "standards", "go", "style.md"), []byte("# Go Style"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := Section{
		Dir: "standards",
		Children: []Section{
			{
				Dir:   "go",
				Files: []File{{Name: "style.md"}},
			},
		},
	}
	if err := validateSectionSources(contentDir, "", s); err != nil {
		t.Fatalf("unexpected error for nested child: %v", err)
	}
}

func TestValidateSectionSources_missingSourceFile(t *testing.T) {
	_, contentDir := makeContentRoot(t)
	s := Section{
		Files: []File{{Name: "nonexistent.md"}},
	}
	if err := validateSectionSources(contentDir, "", s); err == nil {
		t.Fatal("expected error for missing source file; got nil")
	}
}

func TestValidateSectionSources_sourceIsDirectory(t *testing.T) {
	_, contentDir := makeContentRoot(t)
	// create a directory where a file is expected
	if err := os.MkdirAll(filepath.Join(contentDir, "naming.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	s := Section{
		Files: []File{{Name: "naming.md"}},
	}
	if err := validateSectionSources(contentDir, "", s); err == nil {
		t.Fatal("expected error when source path is a directory; got nil")
	}
}

func TestValidateSectionSources_emptyFileName(t *testing.T) {
	_, contentDir := makeContentRoot(t)
	s := Section{
		Files: []File{{Name: ""}},
	}
	if err := validateSectionSources(contentDir, "", s); err == nil {
		t.Fatal("expected error for empty file name; got nil")
	}
}

// ---- ParseFile --------------------------------------------------------

func TestParseFile_validManifest(t *testing.T) {
	tmp := t.TempDir()
	contentDir := filepath.Join(tmp, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "ai_protocol.md"), []byte("# AI Protocol"), 0o644); err != nil {
		t.Fatal(err)
	}

	yamlContent := `version: 1
root_dir: .context
content_root: content
sections:
  - files:
      - name: ai_protocol.md
presets:
  full:
    mode: full
`
	manifestPath := filepath.Join(tmp, "manifest.yml")
	if err := os.WriteFile(manifestPath, []byte(yamlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	m, err := ParseFile(manifestPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Version != 1 {
		t.Errorf("Version = %d; want 1", m.Version)
	}
	if m.RootDir != ".context" {
		t.Errorf("RootDir = %q; want %q", m.RootDir, ".context")
	}
	if _, ok := m.Presets["full"]; !ok {
		t.Error("preset 'full' not found")
	}
	if len(m.Sections) != 1 {
		t.Fatalf("Sections = %d; want 1", len(m.Sections))
	}
	if len(m.Sections[0].Files) != 1 {
		t.Fatalf("Files = %d; want 1", len(m.Sections[0].Files))
	}
	if m.Sections[0].Files[0].Name != "ai_protocol.md" {
		t.Errorf("file name = %q; want %q", m.Sections[0].Files[0].Name, "ai_protocol.md")
	}
}

func TestParseFile_missingFile(t *testing.T) {
	_, err := ParseFile(filepath.Join(t.TempDir(), "nonexistent.yml"))
	if err == nil {
		t.Fatal("expected error for missing manifest file; got nil")
	}
}

func TestParseFile_invalidYAML(t *testing.T) {
	tmp := t.TempDir()
	manifestPath := filepath.Join(tmp, "manifest.yml")
	// structurally invalid YAML — mixing sequence and mapping at the same level
	invalid := "presets:\n  full:\n    mode: full\n  - bad\n"
	if err := os.WriteFile(manifestPath, []byte(invalid), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseFile(manifestPath); err == nil {
		t.Fatal("expected error for invalid YAML; got nil")
	}
}
