package adapter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSupportedText(t *testing.T) {
	got := SupportedText()
	if !strings.Contains(got, AdapterClaude) || !strings.Contains(got, AdapterCodex) {
		t.Fatalf("SupportedText() = %q; want both %q and %q", got, AdapterClaude, AdapterCodex)
	}
	for _, part := range strings.Split(got, ", ") {
		if strings.TrimSpace(part) == "" {
			t.Fatalf("SupportedText() = %q; contains empty entry", got)
		}
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "", wantErr: false},
		{name: AdapterClaude, wantErr: false},
		{name: AdapterCodex, wantErr: false},
		{name: "CLAUDE", wantErr: true},
		{name: "CODEX", wantErr: true},
	}

	for _, tc := range tests {
		err := Validate(tc.name)
		if tc.wantErr && err == nil {
			t.Errorf("Validate(%q) expected error; got nil", tc.name)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("Validate(%q) returned unexpected error: %v", tc.name, err)
		}
	}
}

func TestGenerateEmptyAdapterNameReturnsZeroResult(t *testing.T) {
	res, err := Generate("", t.TempDir(), Options{})
	if err != nil {
		t.Fatalf("Generate returned unexpected error: %v", err)
	}
	if res != (Result{}) {
		t.Fatalf("Generate with empty adapter returned %+v; want zero Result", res)
	}
}

func TestGenerateUnsupportedAdapterReturnsError(t *testing.T) {
	_, err := Generate("unknown", t.TempDir(), Options{})
	if err == nil {
		t.Fatal("Generate expected error for unsupported adapter; got nil")
	}
	if !strings.Contains(err.Error(), "unsupported adapter") {
		t.Fatalf("error = %q; want unsupported adapter message", err)
	}
}

func TestGenerateAdapterFileWritesPrimaryWhenMissing(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := generateAdapterFile("test", root, []byte("hello"), []string{"PRIMARY.md"}, Options{Writer: &out})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}

	wantPath := filepath.Join(root, "PRIMARY.md")
	if res.GeneratedPath != wantPath {
		t.Fatalf("GeneratedPath = %q; want %q", res.GeneratedPath, wantPath)
	}
	if res.UsedFallback {
		t.Fatal("UsedFallback = true; want false")
	}
	if res.Skipped {
		t.Fatal("Skipped = true; want false")
	}
	assertFileContent(t, wantPath, "hello")
}

func TestGenerateAdapterFileUsesFallbackWhenPrimaryExists(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	var out strings.Builder

	res, err := generateAdapterFile("test", root, []byte("generated"), []string{"PRIMARY.md"}, Options{Writer: &out})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}

	wantPath := filepath.Join(root, "PRIMARY.ctx-init.md")
	if res.GeneratedPath != wantPath {
		t.Fatalf("GeneratedPath = %q; want %q", res.GeneratedPath, wantPath)
	}
	if !res.UsedFallback {
		t.Fatal("UsedFallback = false; want true")
	}
	if !strings.Contains(res.Message, "Append or merge") {
		t.Fatalf("Message = %q; want manual merge guidance", res.Message)
	}
	assertFileContent(t, wantPath, "generated")
	assertFileContent(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
}

func TestGenerateAdapterFileSkipsWhenTargetExistsWithoutForce(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.ctx-init.md"), "existing-fallback")
	var out strings.Builder

	res, err := generateAdapterFile("test", root, []byte("generated"), []string{"PRIMARY.md"}, Options{Writer: &out})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if !res.Skipped {
		t.Fatal("Skipped = false; want true")
	}
	if strings.Contains(out.String(), "generated:") {
		t.Fatalf("output = %q; want skipped note without generated label", out.String())
	}
	if !strings.Contains(out.String(), "fallback:") {
		t.Fatalf("output = %q; want fallback label in skipped note", out.String())
	}
	assertFileContent(t, filepath.Join(root, "PRIMARY.ctx-init.md"), "existing-fallback")
}

func TestGenerateAdapterFileOverwritesTargetWithForce(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.ctx-init.md"), "existing-fallback")

	res, err := generateAdapterFile("test", root, []byte("forced"), []string{"PRIMARY.md"}, Options{Force: true})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if !res.UsedFallback {
		t.Fatal("UsedFallback = false; want true")
	}
	if res.Skipped {
		t.Fatal("Skipped = true; want false")
	}
	if res.GeneratedPath != filepath.Join(root, "PRIMARY.ctx-init.md") {
		t.Fatalf("GeneratedPath = %q; want fallback path", res.GeneratedPath)
	}
	assertFileContent(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	assertFileContent(t, filepath.Join(root, "PRIMARY.ctx-init.md"), "forced")
}

func TestGenerateAdapterFileDryRunDoesNotWritePrimary(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder

	res, err := generateAdapterFile("test", root, []byte("hello"), []string{"PRIMARY.md"}, Options{
		DryRun: true,
		Writer: &out,
	})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "PRIMARY.md") {
		t.Fatalf("GeneratedPath = %q; want primary path", res.GeneratedPath)
	}
	if exists, _ := fileExists(filepath.Join(root, "PRIMARY.md")); exists {
		t.Fatal("expected dry run not to create primary file")
	}
}

func TestGenerateAdapterFileDryRunUsesFallbackPathWhenPrimaryExists(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	var out strings.Builder

	res, err := generateAdapterFile("test", root, []byte("hello"), []string{"PRIMARY.md"}, Options{
		DryRun: true,
		Writer: &out,
	})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "PRIMARY.ctx-init.md") {
		t.Fatalf("GeneratedPath = %q; want fallback path", res.GeneratedPath)
	}
	if exists, _ := fileExists(filepath.Join(root, "PRIMARY.ctx-init.md")); exists {
		t.Fatal("expected dry run not to create fallback file")
	}
}

func TestGenerateAdapterFileDryRunSkipsWhenFallbackExistsWithoutForce(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.ctx-init.md"), "existing-fallback")
	var out strings.Builder

	res, err := generateAdapterFile("test", root, []byte("hello"), []string{"PRIMARY.md"}, Options{
		DryRun: true,
		Writer: &out,
	})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if res.Action != ActionDryRunSkip {
		t.Fatalf("Action = %q; want %q", res.Action, ActionDryRunSkip)
	}
	if !res.Skipped {
		t.Fatal("Skipped = false; want true")
	}
	if !strings.Contains(out.String(), "[dry-run/skip]") {
		t.Fatalf("output = %q; want dry-run/skip status", out.String())
	}
	if strings.Contains(out.String(), "generated:") {
		t.Fatalf("output = %q; want dry-run skipped note without generated label", out.String())
	}
	assertFileContent(t, filepath.Join(root, "PRIMARY.ctx-init.md"), "existing-fallback")
}

func TestGenerateAdapterFileDryRunWithForceUsesFallbackPath(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	var out strings.Builder

	res, err := generateAdapterFile("test", root, []byte("hello"), []string{"PRIMARY.md"}, Options{
		DryRun: true,
		Force:  true,
		Writer: &out,
	})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if res.Action != ActionDryRunGenerate {
		t.Fatalf("Action = %q; want %q", res.Action, ActionDryRunGenerate)
	}
	if res.GeneratedPath != filepath.Join(root, "PRIMARY.ctx-init.md") {
		t.Fatalf("GeneratedPath = %q; want fallback path", res.GeneratedPath)
	}
	if !res.UsedFallback {
		t.Fatal("UsedFallback = false; want true")
	}
}

func TestGenerateAdapterFileWithForceWritesPrimaryWhenPrimaryIsMissing(t *testing.T) {
	root := t.TempDir()

	res, err := generateAdapterFile("test", root, []byte("forced"), []string{"PRIMARY.md"}, Options{Force: true})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "PRIMARY.md") {
		t.Fatalf("GeneratedPath = %q; want primary path", res.GeneratedPath)
	}
	if res.UsedFallback {
		t.Fatal("UsedFallback = true; want false")
	}
	assertFileContent(t, filepath.Join(root, "PRIMARY.md"), "forced")
}

func TestGenerateAdapterFileCreatesParentDirectories(t *testing.T) {
	root := t.TempDir()

	res, err := generateAdapterFile("test", root, []byte("nested"), []string{"adapters/PRIMARY.md"}, Options{})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}

	wantPath := filepath.Join(root, "adapters", "PRIMARY.md")
	if res.GeneratedPath != wantPath {
		t.Fatalf("GeneratedPath = %q; want %q", res.GeneratedPath, wantPath)
	}
	assertFileContent(t, wantPath, "nested")
}

func TestGenerateAdapterFileRejectsPathOutsideProjectRoot(t *testing.T) {
	root := t.TempDir()

	_, err := generateAdapterFile("test", root, []byte("bad"), []string{"../PRIMARY.md"}, Options{})
	if err == nil {
		t.Fatal("generateAdapterFile expected error for path outside project root; got nil")
	}
	if !strings.Contains(err.Error(), "test adapter") {
		t.Fatalf("error = %q; want adapter name in validation message", err)
	}
}

func TestGenerateAdapterFileUsesFirstExistingCandidate(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "SECONDARY.md"), "existing-secondary")

	res, err := generateAdapterFile("test", root, []byte("generated"), []string{"PRIMARY.md", "SECONDARY.md"}, Options{})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if res.GeneratedPath != filepath.Join(root, "SECONDARY.ctx-init.md") {
		t.Fatalf("GeneratedPath = %q; want fallback beside first existing candidate", res.GeneratedPath)
	}
}

func TestGenerateAdapterFileRendersAIProtocolPlaceholder(t *testing.T) {
	root := t.TempDir()

	res, err := generateAdapterFile("test", root, []byte("Read {{AI_PROTOCOL_PATH}}"), []string{".claude/PRIMARY.md"}, Options{})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	assertFileContent(t, res.GeneratedPath, "Read ../.context/ai_protocol.md")
}

func writeAdapterFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s content = %q; want %q", path, string(got), want)
	}
}
