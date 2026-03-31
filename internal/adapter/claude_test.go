package adapter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateClaudeCreatesPrimaryFileWhenMissing(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterClaude, root, Options{Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	wantPath := filepath.Join(root, "CLAUDE.md")
	if res.GeneratedPath != wantPath {
		t.Fatalf("GeneratedPath = %q; want %q", res.GeneratedPath, wantPath)
	}
	if res.UsedFallback {
		t.Fatal("UsedFallback = true; want false")
	}
	if !strings.Contains(out.String(), "[generated]") {
		t.Fatalf("output = %q; want generated status", out.String())
	}
}

func TestGenerateClaudeUsesFallbackWhenPrimaryExists(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# Existing\n")
	var out strings.Builder

	res, err := Generate(AdapterClaude, root, Options{Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	wantPath := filepath.Join(root, "CLAUDE.ctx-init.md")
	if res.GeneratedPath != wantPath {
		t.Fatalf("GeneratedPath = %q; want %q", res.GeneratedPath, wantPath)
	}
	if !res.UsedFallback {
		t.Fatal("UsedFallback = false; want true")
	}
	if !strings.Contains(res.Message, "Append or merge") {
		t.Fatalf("Message = %q; want append or merge guidance", res.Message)
	}
	if !strings.Contains(out.String(), "note:") {
		t.Fatalf("output = %q; want note about manual merge", out.String())
	}
}

func TestGenerateClaudeSkipsExistingFallbackWithoutForce(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# Existing\n")
	writeFile(t, filepath.Join(root, "CLAUDE.ctx-init.md"), "# Existing fallback\n")
	var out strings.Builder

	res, err := Generate(AdapterClaude, root, Options{Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if !res.Skipped {
		t.Fatal("Skipped = false; want true")
	}
	if !strings.Contains(out.String(), "[skipped]") {
		t.Fatalf("output = %q; want skipped status", out.String())
	}
}

func TestGenerateClaudeForceDoesNotOverwritePrimary(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# Existing\n")
	writeFile(t, filepath.Join(root, "CLAUDE.ctx-init.md"), "# Existing fallback\n")
	var out strings.Builder

	res, err := Generate(AdapterClaude, root, Options{Force: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "CLAUDE.ctx-init.md") {
		t.Fatalf("GeneratedPath = %q; want CLAUDE.ctx-init.md path", res.GeneratedPath)
	}
	if !res.UsedFallback {
		t.Fatal("UsedFallback = false; want true")
	}
	if !strings.Contains(out.String(), "note:") {
		t.Fatalf("output = %q; want manual merge note", out.String())
	}
	content, err := os.ReadFile(filepath.Join(root, "CLAUDE.ctx-init.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.ctx-init.md: %v", err)
	}
	if !strings.Contains(string(content), "@.context/ai_protocol.md") {
		t.Fatalf("CLAUDE.ctx-init.md content = %q; want adapter template content", string(content))
	}
	assertFileContent(t, filepath.Join(root, "CLAUDE.md"), "# Existing\n")
}

func TestGenerateClaudeWithForceCreatesPrimaryWhenMissing(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterClaude, root, Options{Force: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "CLAUDE.md") {
		t.Fatalf("GeneratedPath = %q; want CLAUDE.md path", res.GeneratedPath)
	}
	if res.UsedFallback {
		t.Fatal("UsedFallback = true; want false")
	}
	if strings.Contains(out.String(), "note:") {
		t.Fatalf("output = %q; want no manual merge note", out.String())
	}
	content, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	if !strings.Contains(string(content), "@.context/ai_protocol.md") {
		t.Fatalf("CLAUDE.md content = %q; want adapter template content", string(content))
	}
}

func TestGenerateClaudeDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterClaude, root, Options{DryRun: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "CLAUDE.md") {
		t.Fatalf("GeneratedPath = %q; want CLAUDE.md path", res.GeneratedPath)
	}
	if exists, _ := fileExists(filepath.Join(root, "CLAUDE.md")); exists {
		t.Fatal("expected no file to be written during dry run")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
