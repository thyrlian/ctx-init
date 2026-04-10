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

	wantPath := filepath.Join(root, ".claude", "CLAUDE.md")
	if res.GeneratedPath != wantPath {
		t.Fatalf("GeneratedPath = %q; want %q", res.GeneratedPath, wantPath)
	}
	if res.UsedFallback {
		t.Fatal("UsedFallback = true; want false")
	}
	if !strings.Contains(out.String(), "[generated]") {
		t.Fatalf("output = %q; want generated status", out.String())
	}
	content, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read generated CLAUDE.md: %v", err)
	}
	assertClaudeTemplateContent(t, string(content), "@../.context/ai_protocol.md", "@../../../../.context/ai_protocol.md")
}

func TestGenerateClaudeUsesFallbackWhenRootPrimaryExists(t *testing.T) {
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
	content, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read generated fallback CLAUDE.ctx-init.md: %v", err)
	}
	assertClaudeTemplateContent(t, string(content), "@.context/ai_protocol.md", "@../../../.context/ai_protocol.md")
}

func TestGenerateClaudePrefersScopedPrimaryWhenBothLocationsExist(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".claude", "CLAUDE.md"), "# Existing scoped\n")
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# Existing root\n")

	res, err := Generate(AdapterClaude, root, Options{})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, ".claude", "CLAUDE.ctx-init.md") {
		t.Fatalf("GeneratedPath = %q; want scoped fallback path", res.GeneratedPath)
	}
	assertFileContent(t, filepath.Join(root, ".claude", "CLAUDE.md"), "# Existing scoped\n")
	assertFileContent(t, filepath.Join(root, "CLAUDE.md"), "# Existing root\n")
}

func TestGenerateClaudeWithForceCreatesPrimaryWhenMissing(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterClaude, root, Options{Force: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, ".claude", "CLAUDE.md") {
		t.Fatalf("GeneratedPath = %q; want CLAUDE.md path", res.GeneratedPath)
	}
	if res.UsedFallback {
		t.Fatal("UsedFallback = true; want false")
	}
	if strings.Contains(out.String(), "note:") {
		t.Fatalf("output = %q; want no manual merge note", out.String())
	}
	content, err := os.ReadFile(filepath.Join(root, ".claude", "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	assertClaudeTemplateContent(t, string(content), "@../.context/ai_protocol.md", "@../../../../.context/ai_protocol.md")
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func assertClaudeTemplateContent(t *testing.T, content, normalPath, worktreePath string) {
	t.Helper()
	if !strings.Contains(content, normalPath) {
		t.Fatalf("CLAUDE content = %q; want normal ai_protocol path %q", content, normalPath)
	}
	if !strings.Contains(content, worktreePath) {
		t.Fatalf("CLAUDE content = %q; want worktree ai_protocol path %q", content, worktreePath)
	}
	if !strings.Contains(content, ".claude/worktrees/<worktree>/.claude/") {
		t.Fatalf("CLAUDE content = %q; want worktree location hint", content)
	}
	if strings.Contains(content, "{{AI_PROTOCOL_PATH}}") {
		t.Fatalf("CLAUDE content = %q; want AI_PROTOCOL_PATH placeholder to be fully rendered", content)
	}
}
