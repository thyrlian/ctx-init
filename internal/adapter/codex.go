package adapter

import (
	"fmt"

	assetdata "github.com/thyrlian/ctx-init/assets"
)

const (
	codexTemplateName = "AGENTS.md"
)

var codexCandidates = []string{
	"AGENTS.md",
}

func generateCodex(projectRoot string, opt Options) (Result, error) {
	content, err := assetdata.ReadAdapter(codexTemplateName)
	if err != nil {
		return Result{}, fmt.Errorf("read Codex adapter template: %w", err)
	}

	return generateAdapterFile(AdapterCodex, projectRoot, content, codexCandidates, opt)
}
