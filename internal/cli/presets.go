package cli

import (
	"fmt"
	"strings"
)

const (
	PresetMinimal  = "minimal"
	PresetStandard = "standard"
	PresetFull     = "full"
)

const DefaultPreset = PresetStandard

var supportedPresets = []string{
	PresetMinimal,
	PresetStandard,
	PresetFull,
}

func supportedPresetsText() string {
	return strings.Join(supportedPresets, ", ")
}

func validatePreset(preset string) error {
	for _, p := range supportedPresets {
		if preset == p {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid preset %q (supported: %s)",
		preset,
		supportedPresetsText(),
	)
}
