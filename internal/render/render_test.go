package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/thyrlian/ctx-init/internal/manifest"
	"github.com/thyrlian/ctx-init/internal/plan"
)

// ---- appendCtxID ------------------------------------------------------

func TestAppendCtxID_appendsToEmptyFile(t *testing.T) {
	f := tempFile(t, "")
	token, err := appendCtxID(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token) != 16 {
		t.Errorf("token length = %d; want 16", len(token))
	}
	content := readFile(t, f)
	if !strings.Contains(content, "<!-- ctx-id: "+token+" -->") {
		t.Errorf("token not found in file content:\n%s", content)
	}
}

func TestAppendCtxID_appendsToExistingContent(t *testing.T) {
	f := tempFile(t, "# Hello\n\nSome content.\n")
	token, err := appendCtxID(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readFile(t, f)
	if !strings.Contains(content, "# Hello") {
		t.Error("original content was lost")
	}
	if !strings.Contains(content, "<!-- ctx-id: "+token+" -->") {
		t.Errorf("ctx-id not found:\n%s", content)
	}
}

func TestAppendCtxID_replacesExistingToken(t *testing.T) {
	f := tempFile(t, "# Hello\n\n<!-- ctx-id: aabbccdd11223344 -->\n")
	token, err := appendCtxID(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "aabbccdd11223344" {
		t.Fatal("new token should differ from original (collision is astronomically unlikely)")
	}
	content := readFile(t, f)
	if strings.Contains(content, "aabbccdd11223344") {
		t.Error("old ctx-id was not replaced")
	}
	if strings.Count(content, "<!-- ctx-id:") != 1 {
		t.Errorf("expected exactly one ctx-id, got:\n%s", content)
	}
}

func TestAppendCtxID_tokenIs16HexChars(t *testing.T) {
	f := tempFile(t, "content\n")
	token, err := appendCtxID(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token) != 16 {
		t.Errorf("token length = %d; want 16", len(token))
	}
	for _, c := range token {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Errorf("token %q contains non-hex character %q", token, c)
		}
	}
}

func TestAppendCtxID_noTrailingNewlineInOriginal(t *testing.T) {
	// File without trailing newline — ensure ctx-id is still on its own line
	f := tempFile(t, "# No newline at end")
	token, err := appendCtxID(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readFile(t, f)
	expected := "<!-- ctx-id: " + token + " -->"
	if !strings.Contains(content, "\n\n"+expected+"\n") {
		t.Errorf("ctx-id not preceded by blank line or not on its own line:\n%s", content)
	}
}

// ---- buildIndexContent ------------------------------------------------

func TestBuildIndexContent_header(t *testing.T) {
	p := &plan.Plan{
		Preset:     "standard",
		OutRootAbs: "/tmp/out",
		Entries:    nil,
	}
	content, err := buildIndexContent(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(content, "# Context Index\n") {
		t.Errorf("missing h1 header:\n%s", content)
	}
	if !strings.Contains(content, "preset: `standard`") {
		t.Errorf("preset not found in header:\n%s", content)
	}
	if !strings.Contains(content, "## Files\n") {
		t.Errorf("missing Files section:\n%s", content)
	}
}

func TestBuildIndexContent_fileListWithTags(t *testing.T) {
	outRoot := t.TempDir()
	p := &plan.Plan{
		Preset:     "full",
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Dst:     filepath.Join(outRoot, "ai_protocol.md"),
				Tags:    []string{"core"},
				Gateway: false,
			},
			{
				Dst:     filepath.Join(outRoot, "standards", "naming.md"),
				Tags:    []string{"standards"},
				Gateway: false,
			},
		},
	}
	content, err := buildIndexContent(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(content, "ai_protocol.md") {
		t.Errorf("ai_protocol.md not found in index:\n%s", content)
	}
	if !strings.Contains(content, "`core`") {
		t.Errorf("`core` tag not found:\n%s", content)
	}
	if !strings.Contains(content, "standards/naming.md") {
		t.Errorf("standards/naming.md not found:\n%s", content)
	}
	if !strings.Contains(content, "`standards`") {
		t.Errorf("`standards` tag not found:\n%s", content)
	}
}

func TestBuildIndexContent_gatewayTagAppended(t *testing.T) {
	outRoot := t.TempDir()
	p := &plan.Plan{
		Preset:     "full",
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Dst:     filepath.Join(outRoot, "standards", "styleguide.md"),
				Tags:    []string{"standards"},
				Gateway: true,
			},
		},
	}
	content, err := buildIndexContent(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(content, "`gateway`") {
		t.Errorf("`gateway` tag not appended for gateway entry:\n%s", content)
	}
}

func TestBuildIndexContent_noTagsEntry(t *testing.T) {
	outRoot := t.TempDir()
	p := &plan.Plan{
		Preset:     "full",
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Dst:     filepath.Join(outRoot, "notags.md"),
				Tags:    nil,
				Gateway: false,
			},
		},
	}
	content, err := buildIndexContent(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Entry with no tags should still appear as a link without trailing tags
	if !strings.Contains(content, "[notags.md](notags.md)") {
		t.Errorf("tagless entry not found or malformed:\n%s", content)
	}
}

// ---- fileExists -------------------------------------------------------

func TestFileExists_existingFile(t *testing.T) {
	f := tempFile(t, "data")
	ok, err := fileExists(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true for existing file")
	}
}

func TestFileExists_missingFile(t *testing.T) {
	ok, err := fileExists(filepath.Join(t.TempDir(), "does_not_exist.md"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false for missing file")
	}
}

// ---- Action.String ----------------------------------------------------

func TestActionString(t *testing.T) {
	tests := []struct {
		action Action
		want   string
	}{
		{ActionGenerated, "generated"},
		{ActionCopied, "copied"},
		{ActionSkipped, "skipped"},
		{ActionDryRunGenerate, "dry-run/generate"},
		{ActionDryRunCopy, "dry-run/copy"},
		{ActionDryRunSkip, "dry-run/skip"},
		{ActionUnknown, "unknown"},
	}
	for _, tc := range tests {
		if got := tc.action.String(); got != tc.want {
			t.Errorf("Action(%d).String() = %q; want %q", tc.action, got, tc.want)
		}
	}
}

// ---- Run (integration) ------------------------------------------------

func TestRun_copiesFilesAndGeneratesIndex(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()
	outRoot := filepath.Join(outDir, ".context")

	// Create source files
	writeFile(t, filepath.Join(srcDir, "ai_protocol.md"), "# AI Protocol\n")
	writeFile(t, filepath.Join(srcDir, "standards", "naming.md"), "# Naming\n")

	p := &plan.Plan{
		Preset:     "full",
		Mode:       manifest.ModeFull,
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Src:  filepath.Join(srcDir, "ai_protocol.md"),
				Dst:  filepath.Join(outRoot, "ai_protocol.md"),
				Tags: []string{"core"},
			},
			{
				Src:  filepath.Join(srcDir, "standards", "naming.md"),
				Dst:  filepath.Join(outRoot, "standards", "naming.md"),
				Tags: []string{"standards"},
			},
		},
	}

	var buf strings.Builder
	res, err := Run(p, Options{Writer: &buf})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if res.Copied != 2 {
		t.Errorf("Copied = %d; want 2", res.Copied)
	}
	if res.Generated != 1 {
		t.Errorf("Generated = %d; want 1", res.Generated)
	}
	if res.Skipped != 0 {
		t.Errorf("Skipped = %d; want 0", res.Skipped)
	}

	// ctx-id should be recorded for .md files
	if len(res.FileIDs) != 3 {
		t.Errorf("FileIDs has %d entries; want 3: %v", len(res.FileIDs), res.FileIDs)
	}

	// recorded ctx-id must match what was actually written into each file
	for relPath, id := range res.FileIDs {
		dstPath := filepath.Join(outRoot, filepath.FromSlash(relPath))
		content := readFile(t, dstPath)
		expected := "<!-- ctx-id: " + id + " -->"
		if !strings.Contains(content, expected) {
			t.Errorf("%s: file content missing ctx-id %q:\n%s", relPath, id, content)
		}
	}

	// output must mention both the copy and generate actions
	if !strings.Contains(buf.String(), "[copied]") {
		t.Errorf("output missing [copied]:\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "[generated]") {
		t.Errorf("output missing [generated]:\n%s", buf.String())
	}

	// _INDEX.md must exist
	indexPath := filepath.Join(outRoot, "_INDEX.md")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("_INDEX.md was not generated")
	}
	if indexID := res.FileIDs["_INDEX.md"]; indexID == "" {
		t.Fatalf("_INDEX.md ctx-id was not recorded: %v", res.FileIDs)
	} else {
		content := readFile(t, indexPath)
		expected := "<!-- ctx-id: " + indexID + " -->"
		if !strings.Contains(content, expected) {
			t.Errorf("_INDEX.md missing ctx-id %q:\n%s", indexID, content)
		}
	}
}

func TestRun_dryRunDoesNotWriteFiles(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()
	outRoot := filepath.Join(outDir, ".context")

	writeFile(t, filepath.Join(srcDir, "ai_protocol.md"), "# AI Protocol\n")

	p := &plan.Plan{
		Preset:     "full",
		Mode:       manifest.ModeFull,
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Src:  filepath.Join(srcDir, "ai_protocol.md"),
				Dst:  filepath.Join(outRoot, "ai_protocol.md"),
				Tags: []string{"core"},
			},
		},
	}

	var buf strings.Builder
	_, err := Run(p, Options{DryRun: true, Writer: &buf})
	if err != nil {
		t.Fatalf("Run dry-run failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outRoot, "ai_protocol.md")); !os.IsNotExist(err) {
		t.Error("dry-run must not write files to disk")
	}
	if _, err := os.Stat(filepath.Join(outRoot, "_INDEX.md")); !os.IsNotExist(err) {
		t.Error("dry-run must not write _INDEX.md to disk")
	}
	if !strings.Contains(buf.String(), "[dry-run/copy]") {
		t.Errorf("output missing [dry-run/copy]:\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "[dry-run/generate]") {
		t.Errorf("output missing [dry-run/generate]:\n%s", buf.String())
	}
}

func TestRun_skipsExistingWithoutForce(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()
	outRoot := filepath.Join(outDir, ".context")
	if err := os.MkdirAll(outRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	writeFile(t, filepath.Join(srcDir, "ai_protocol.md"), "# Source\n")
	writeFile(t, filepath.Join(outRoot, "ai_protocol.md"), "# Existing\n")

	p := &plan.Plan{
		Preset:     "full",
		Mode:       manifest.ModeFull,
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Src:  filepath.Join(srcDir, "ai_protocol.md"),
				Dst:  filepath.Join(outRoot, "ai_protocol.md"),
				Tags: []string{"core"},
			},
		},
	}

	var buf strings.Builder
	res, err := Run(p, Options{Writer: &buf})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("Skipped = %d; want 1", res.Skipped)
	}
	if !strings.Contains(buf.String(), "[skipped]") {
		t.Errorf("output missing [skipped]:\n%s", buf.String())
	}
	// Content should be unchanged
	if content := readFile(t, filepath.Join(outRoot, "ai_protocol.md")); !strings.Contains(content, "# Existing") {
		t.Error("existing file was overwritten without --force")
	}
}

func TestRun_forceOverwritesExisting(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()
	outRoot := filepath.Join(outDir, ".context")
	if err := os.MkdirAll(outRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	writeFile(t, filepath.Join(srcDir, "ai_protocol.md"), "# New content\n")
	writeFile(t, filepath.Join(outRoot, "ai_protocol.md"), "# Old content\n")

	p := &plan.Plan{
		Preset:     "full",
		Mode:       manifest.ModeFull,
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Src:  filepath.Join(srcDir, "ai_protocol.md"),
				Dst:  filepath.Join(outRoot, "ai_protocol.md"),
				Tags: []string{"core"},
			},
		},
	}

	var buf strings.Builder
	res, err := Run(p, Options{Force: true, Writer: &buf})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if res.Copied != 1 {
		t.Errorf("Copied = %d; want 1", res.Copied)
	}
	content := readFile(t, filepath.Join(outRoot, "ai_protocol.md"))
	if !strings.Contains(content, "# New content") {
		t.Errorf("force did not overwrite; got:\n%s", content)
	}
}

func TestRun_nonMarkdownFileCopiedWithoutCtxID(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()
	outRoot := filepath.Join(outDir, ".context")

	writeFile(t, filepath.Join(srcDir, "config.json"), `{"ok":true}`)

	p := &plan.Plan{
		Preset:     "full",
		Mode:       manifest.ModeFull,
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Src: filepath.Join(srcDir, "config.json"),
				Dst: filepath.Join(outRoot, "config.json"),
			},
		},
	}

	var buf strings.Builder
	res, err := Run(p, Options{Writer: &buf})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if res.Copied != 1 {
		t.Errorf("Copied = %d; want 1", res.Copied)
	}
	if len(res.FileIDs) != 1 {
		t.Errorf("FileIDs = %v; want only _INDEX.md to get a ctx-id", res.FileIDs)
	}
	if _, ok := res.FileIDs["_INDEX.md"]; !ok {
		t.Errorf("FileIDs = %v; want _INDEX.md ctx-id to be recorded", res.FileIDs)
	}
	if _, ok := res.FileIDs["config.json"]; ok {
		t.Errorf("FileIDs = %v; non-markdown file must not get ctx-id", res.FileIDs)
	}
	if !strings.Contains(buf.String(), "[copied]") {
		t.Errorf("output missing [copied]:\n%s", buf.String())
	}
	content := readFile(t, filepath.Join(outRoot, "config.json"))
	if strings.Contains(content, "ctx-id") {
		t.Errorf("non-markdown file must not contain ctx-id:\n%s", content)
	}
}

func TestRun_dryRunSkipsExistingWithoutForce(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()
	outRoot := filepath.Join(outDir, ".context")
	if err := os.MkdirAll(outRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	writeFile(t, filepath.Join(srcDir, "ai_protocol.md"), "# Source\n")
	writeFile(t, filepath.Join(outRoot, "ai_protocol.md"), "# Existing\n")

	p := &plan.Plan{
		Preset:     "full",
		Mode:       manifest.ModeFull,
		OutRootAbs: outRoot,
		Entries: []plan.Entry{
			{
				Src: filepath.Join(srcDir, "ai_protocol.md"),
				Dst: filepath.Join(outRoot, "ai_protocol.md"),
			},
		},
	}

	var buf strings.Builder
	res, err := Run(p, Options{DryRun: true, Writer: &buf})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("Skipped = %d; want 1", res.Skipped)
	}
	if !strings.Contains(buf.String(), "[dry-run/skip]") {
		t.Errorf("output missing [dry-run/skip]:\n%s", buf.String())
	}
	// Existing file must be untouched
	content := readFile(t, filepath.Join(outRoot, "ai_protocol.md"))
	if !strings.Contains(content, "# Existing") {
		t.Error("dry-run/skip must not modify existing file")
	}
}

// ---- helpers ----------------------------------------------------------

func tempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "ctx-test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	if content != "" {
		if _, err := f.WriteString(content); err != nil {
			t.Fatal(err)
		}
	}
	f.Close()
	return f.Name()
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
