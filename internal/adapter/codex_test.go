package adapter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCodexCreatesPrimaryFileWhenMissing(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterCodex, root, Options{Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	wantPath := filepath.Join(root, "AGENTS.md")
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

func TestGenerateCodexUsesFallbackWhenPrimaryExists(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "AGENTS.md"), "# Existing\n")
	var out strings.Builder

	res, err := Generate(AdapterCodex, root, Options{Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	wantPath := filepath.Join(root, "AGENTS.ctx-init.md")
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

func TestGenerateCodexSkipsExistingFallbackWithoutForce(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "AGENTS.md"), "# Existing\n")
	writeFile(t, filepath.Join(root, "AGENTS.ctx-init.md"), "# Existing fallback\n")
	var out strings.Builder

	res, err := Generate(AdapterCodex, root, Options{Writer: &out})
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

func TestGenerateCodexForceDoesNotOverwritePrimary(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "AGENTS.md"), "# Existing\n")
	writeFile(t, filepath.Join(root, "AGENTS.ctx-init.md"), "# Existing fallback\n")
	var out strings.Builder

	res, err := Generate(AdapterCodex, root, Options{Force: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "AGENTS.ctx-init.md") {
		t.Fatalf("GeneratedPath = %q; want AGENTS.ctx-init.md path", res.GeneratedPath)
	}
	if !res.UsedFallback {
		t.Fatal("UsedFallback = false; want true")
	}
	if !strings.Contains(out.String(), "note:") {
		t.Fatalf("output = %q; want manual merge note", out.String())
	}
	content, err := os.ReadFile(filepath.Join(root, "AGENTS.ctx-init.md"))
	if err != nil {
		t.Fatalf("read AGENTS.ctx-init.md: %v", err)
	}
	if !strings.Contains(string(content), "Read `.context/ai_protocol.md` before doing any work.") {
		t.Fatalf("AGENTS.ctx-init.md content = %q; want adapter template content", string(content))
	}
	assertFileContent(t, filepath.Join(root, "AGENTS.md"), "# Existing\n")
}

func TestGenerateCodexWithForceCreatesPrimaryWhenMissing(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterCodex, root, Options{Force: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "AGENTS.md") {
		t.Fatalf("GeneratedPath = %q; want AGENTS.md path", res.GeneratedPath)
	}
	if res.UsedFallback {
		t.Fatal("UsedFallback = true; want false")
	}
	if strings.Contains(out.String(), "note:") {
		t.Fatalf("output = %q; want no manual merge note", out.String())
	}
	content, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	if !strings.Contains(string(content), "Read `.context/ai_protocol.md` before doing any work.") {
		t.Fatalf("AGENTS.md content = %q; want adapter template content", string(content))
	}
}

func TestGenerateCodexDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterCodex, root, Options{DryRun: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "AGENTS.md") {
		t.Fatalf("GeneratedPath = %q; want AGENTS.md path", res.GeneratedPath)
	}
	if exists, _ := fileExists(filepath.Join(root, "AGENTS.md")); exists {
		t.Fatal("expected no file to be written during dry run")
	}
}
