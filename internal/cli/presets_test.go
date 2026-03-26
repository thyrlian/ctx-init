package cli

import "testing"

func TestValidatePreset_validPresets(t *testing.T) {
	for _, p := range []string{PresetMinimal, PresetStandard, PresetFull} {
		if err := validatePreset(p); err != nil {
			t.Errorf("validatePreset(%q) returned unexpected error: %v", p, err)
		}
	}
}

func TestValidatePreset_invalidPresets(t *testing.T) {
	invalid := []string{"", "STANDARD", "Standard", " standard ", "all", "none"}
	for _, p := range invalid {
		if err := validatePreset(p); err == nil {
			t.Errorf("validatePreset(%q) expected error; got nil", p)
		}
	}
}
