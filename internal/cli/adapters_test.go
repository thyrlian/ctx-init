package cli

import "testing"

func TestValidateAdapterValidValues(t *testing.T) {
	for _, name := range []string{"", "claude", "codex", "gemini"} {
		if err := validateAdapter(name); err != nil {
			t.Errorf("validateAdapter(%q) returned unexpected error: %v", name, err)
		}
	}
}

func TestValidateAdapterInvalidValues(t *testing.T) {
	for _, name := range []string{"CLAUDE", "CODEX", "GEMINI", " Claude ", "all"} {
		if err := validateAdapter(name); err == nil {
			t.Errorf("validateAdapter(%q) expected error; got nil", name)
		}
	}
}
