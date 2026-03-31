package adapter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSupportedText(t *testing.T) {
	if got := SupportedText(); got != AdapterClaude {
		t.Fatalf("SupportedText() = %q; want %q", got, AdapterClaude)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "", wantErr: false},
		{name: AdapterClaude, wantErr: false},
		{name: "CLAUDE", wantErr: true},
		{name: "codex", wantErr: true},
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

	res, err := generateAdapterFile(root, []byte("hello"), "PRIMARY.md", "PRIMARY.ctx-init.md", Options{Writer: &out})
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

	res, err := generateAdapterFile(root, []byte("generated"), "PRIMARY.md", "PRIMARY.ctx-init.md", Options{Writer: &out})
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

	res, err := generateAdapterFile(root, []byte("generated"), "PRIMARY.md", "PRIMARY.ctx-init.md", Options{Writer: &out})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
	}
	if !res.Skipped {
		t.Fatal("Skipped = false; want true")
	}
	assertFileContent(t, filepath.Join(root, "PRIMARY.ctx-init.md"), "existing-fallback")
}

func TestGenerateAdapterFileOverwritesTargetWithForce(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.ctx-init.md"), "existing-fallback")

	res, err := generateAdapterFile(root, []byte("forced"), "PRIMARY.md", "PRIMARY.ctx-init.md", Options{Force: true})
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

	res, err := generateAdapterFile(root, []byte("hello"), "PRIMARY.md", "PRIMARY.ctx-init.md", Options{
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

	res, err := generateAdapterFile(root, []byte("hello"), "PRIMARY.md", "PRIMARY.ctx-init.md", Options{
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

func TestGenerateAdapterFileDryRunWithForceUsesFallbackPath(t *testing.T) {
	root := t.TempDir()
	writeAdapterFile(t, filepath.Join(root, "PRIMARY.md"), "existing-primary")
	var out strings.Builder

	res, err := generateAdapterFile(root, []byte("hello"), "PRIMARY.md", "PRIMARY.ctx-init.md", Options{
		DryRun: true,
		Force:  true,
		Writer: &out,
	})
	if err != nil {
		t.Fatalf("generateAdapterFile returned error: %v", err)
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

	res, err := generateAdapterFile(root, []byte("forced"), "PRIMARY.md", "PRIMARY.ctx-init.md", Options{Force: true})
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

	res, err := generateAdapterFile(root, []byte("nested"), "adapters/PRIMARY.md", "adapters/PRIMARY.ctx-init.md", Options{})
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

	_, err := generateAdapterFile(root, []byte("bad"), "../PRIMARY.md", "PRIMARY.ctx-init.md", Options{})
	if err == nil {
		t.Fatal("generateAdapterFile expected error for path outside project root; got nil")
	}
	if !strings.Contains(err.Error(), "project root") {
		t.Fatalf("error = %q; want project root validation message", err)
	}
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
