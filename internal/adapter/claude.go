package adapter

import (
	"fmt"

	assetdata "github.com/thyrlian/ctx-init/assets"
)

const (
	claudeTemplateName = "CLAUDE.md"
	claudePrimaryName  = "CLAUDE.md"
	claudeFallbackName = "CLAUDE.ctx-init.md"
)

func generateClaude(projectRoot string, opt Options) (Result, error) {
	content, err := assetdata.ReadAdapter(claudeTemplateName)
	if err != nil {
		return Result{}, fmt.Errorf("read Claude adapter template: %w", err)
	}

	return generateAdapterFile(projectRoot, content, claudePrimaryName, claudeFallbackName, opt)
}
