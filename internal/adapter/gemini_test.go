package adapter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateGeminiCreatesPrimaryFileWhenMissing(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterGemini, root, Options{Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	wantPath := filepath.Join(root, ".agents", "rules", "GEMINI.md")
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
		t.Fatalf("read GEMINI.md: %v", err)
	}
	if !strings.Contains(string(content), "activation: always") {
		t.Fatalf("GEMINI.md content = %q; want Always On activation", string(content))
	}
	if !strings.Contains(string(content), "../../.context/ai_protocol.md") {
		t.Fatalf("GEMINI.md content = %q; want relative ai_protocol path for .agents/rules location", string(content))
	}
}

func TestGenerateGeminiUsesFallbackWhenPrimaryExists(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".agents", "rules", "GEMINI.md"), "# Existing\n")
	var out strings.Builder

	res, err := Generate(AdapterGemini, root, Options{Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	wantPath := filepath.Join(root, ".agents", "rules", "GEMINI.ctx-init.md")
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
	content, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read GEMINI.ctx-init.md: %v", err)
	}
	if !strings.Contains(string(content), "../../.context/ai_protocol.md") {
		t.Fatalf("GEMINI.ctx-init.md content = %q; want relative ai_protocol path for .agents/rules location", string(content))
	}
}

func TestGenerateGeminiSkipsExistingFallbackWithoutForce(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".agents", "rules", "GEMINI.md"), "# Existing\n")
	writeFile(t, filepath.Join(root, ".agents", "rules", "GEMINI.ctx-init.md"), "# Existing fallback\n")
	var out strings.Builder

	res, err := Generate(AdapterGemini, root, Options{Writer: &out})
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

func TestGenerateGeminiForceDoesNotOverwritePrimary(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".agents", "rules", "GEMINI.md"), "# Existing\n")
	writeFile(t, filepath.Join(root, ".agents", "rules", "GEMINI.ctx-init.md"), "# Existing fallback\n")
	var out strings.Builder

	res, err := Generate(AdapterGemini, root, Options{Force: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, ".agents", "rules", "GEMINI.ctx-init.md") {
		t.Fatalf("GeneratedPath = %q; want GEMINI.ctx-init.md path", res.GeneratedPath)
	}
	if !res.UsedFallback {
		t.Fatal("UsedFallback = false; want true")
	}
	if !strings.Contains(out.String(), "note:") {
		t.Fatalf("output = %q; want manual merge note", out.String())
	}
	assertFileContent(t, filepath.Join(root, ".agents", "rules", "GEMINI.md"), "# Existing\n")
}

func TestGenerateGeminiWithForceCreatesPrimaryWhenMissing(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterGemini, root, Options{Force: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, ".agents", "rules", "GEMINI.md") {
		t.Fatalf("GeneratedPath = %q; want GEMINI.md path", res.GeneratedPath)
	}
	if res.UsedFallback {
		t.Fatal("UsedFallback = true; want false")
	}
	if strings.Contains(out.String(), "note:") {
		t.Fatalf("output = %q; want no manual merge note", out.String())
	}
}

func TestGenerateGeminiDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := Generate(AdapterGemini, root, Options{DryRun: true, Writer: &out})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, ".agents", "rules", "GEMINI.md") {
		t.Fatalf("GeneratedPath = %q; want GEMINI.md path", res.GeneratedPath)
	}
	if exists, _ := fileExists(filepath.Join(root, ".agents", "rules", "GEMINI.md")); exists {
		t.Fatal("expected no file to be written during dry run")
	}
}
