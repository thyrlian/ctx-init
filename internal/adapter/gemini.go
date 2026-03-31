package adapter

import (
	"fmt"

	assetdata "github.com/thyrlian/ctx-init/assets"
)

const (
	geminiTemplateName = "GEMINI.md"
)

var geminiCandidates = []string{
	".agents/rules/GEMINI.md",
}

func generateGemini(projectRoot string, opt Options) (Result, error) {
	content, err := assetdata.ReadAdapter(geminiTemplateName)
	if err != nil {
		return Result{}, fmt.Errorf("read Gemini adapter template: %w", err)
	}

	return generateAdapterFile(AdapterGemini, projectRoot, content, geminiCandidates, opt)
}
