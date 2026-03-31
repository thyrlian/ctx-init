package cli

import "github.com/thyrlian/ctx-init/internal/adapter"

const DefaultAdapter = ""

func supportedAdaptersText() string {
	return adapter.SupportedText()
}

func validateAdapter(name string) error {
	return adapter.Validate(name)
}
