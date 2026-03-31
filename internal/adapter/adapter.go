package adapter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	AdapterClaude = "claude"
	AdapterCodex  = "codex"
)

type generator func(projectRoot string, opt Options) (Result, error)

const aiProtocolPathPlaceholder = "{{AI_PROTOCOL_PATH}}"

var supportedAdapters = []string{
	AdapterClaude,
	AdapterCodex,
}

var generators = map[string]generator{
	AdapterClaude: generateClaude,
	AdapterCodex:  generateCodex,
}

func SupportedText() string {
	return strings.Join(supportedAdapters, ", ")
}

func Validate(name string) error {
	if name == "" {
		return nil
	}

	for _, supported := range supportedAdapters {
		if name == supported {
			return nil
		}
	}

	return fmt.Errorf("invalid adapter %q (supported: %s)", name, SupportedText())
}

func Generate(name, projectRoot string, opt Options) (Result, error) {
	if name == "" {
		return Result{}, nil
	}

	generate, ok := generators[name]
	if !ok {
		return Result{}, fmt.Errorf("unsupported adapter %q", name)
	}

	return generate(projectRoot, opt)
}

func writerOrStdout(w io.Writer) io.Writer {
	if w != nil {
		return w
	}
	return os.Stdout
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func absProjectRoot(projectRoot string) (string, error) {
	if strings.TrimSpace(projectRoot) == "" {
		return "", fmt.Errorf("project root must not be empty")
	}

	rootAbs, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("resolve project root %q: %w", projectRoot, err)
	}
	return rootAbs, nil
}

func projectRelativePath(name string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(strings.TrimSpace(name)))
	if clean == "." || clean == "" {
		return "", fmt.Errorf("path must not be empty")
	}
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("path %q must be relative to the project root", name)
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q must stay within the project root", name)
	}
	return clean, nil
}

func fallbackPathFor(primaryRel string) string {
	ext := filepath.Ext(primaryRel)
	if ext == "" {
		return primaryRel + ".ctx-init"
	}
	base := strings.TrimSuffix(primaryRel, ext)
	return base + ".ctx-init" + ext
}

func aiProtocolPathFor(targetRel string) string {
	rel, err := filepath.Rel(filepath.Dir(targetRel), filepath.FromSlash(".context/ai_protocol.md"))
	if err != nil {
		return filepath.ToSlash(".context/ai_protocol.md")
	}
	return filepath.ToSlash(rel)
}

func renderAdapterTemplate(content []byte, primaryRel string) []byte {
	return []byte(strings.ReplaceAll(string(content), aiProtocolPathPlaceholder, aiProtocolPathFor(primaryRel)))
}

func resolvePrimaryLocation(adapterName string, rootAbs string, candidates []string) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("%s adapter must define at least one candidate location", adapterName)
	}

	// Candidates are adapter-defined paths. Validate them up front so future
	// adapter additions fail fast if they declare an invalid project-relative
	// location.
	validated := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		rel, err := projectRelativePath(candidate)
		if err != nil {
			return "", fmt.Errorf("invalid %s adapter candidate path %q: %w", adapterName, candidate, err)
		}
		validated = append(validated, rel)
	}

	for _, candidate := range validated {
		exists, err := fileExists(filepath.Join(rootAbs, candidate))
		if err != nil {
			return "", fmt.Errorf("check %s adapter candidate %s: %w", adapterName, filepath.Join(rootAbs, candidate), err)
		}
		if exists {
			// Prefer the first official location the user is already using. Only
			// fall back to the first candidate when none of the known locations
			// exist yet.
			return candidate, nil
		}
	}

	return validated[0], nil
}

func printAdapterNote(w io.Writer, action Action, primaryPath, targetPath string) {
	fmt.Fprintf(w, "  note:\n")
	fmt.Fprintf(w, "    existing:  %s\n", primaryPath)
	switch action {
	case ActionGenerated, ActionDryRunGenerate:
		fmt.Fprintf(w, "    generated: %s\n", targetPath)
		fmt.Fprintf(w, "    next:      append or merge the generated file into your existing AI agent instructions file.\n")
	case ActionSkipped, ActionDryRunSkip:
		fmt.Fprintf(w, "    fallback:  %s\n", targetPath)
		fmt.Fprintf(w, "    next:      reuse that fallback file or rerun with -force to replace it before merging into your existing AI agent instructions file.\n")
	default:
		fmt.Fprintf(w, "    target:    %s\n", targetPath)
	}
}

// generateAdapterFile applies the standard ctx-init adapter behavior for a
// single output file:
//   - write the primary file when it does not exist
//   - fall back to an alternate filename when the primary file already exists
//   - respect dry-run semantics and only use force to overwrite an existing
//     fallback file; an existing primary file is never overwritten implicitly
//   - print a user-facing note when manual merge is required
//
// This keeps the common adapter location, rendering, and fallback flow in one
// place so individual adapters only need to supply their content and candidate
// locations. Candidate paths must be project-root-relative file paths.
//
// Terminology:
//   - primary: the selected output filename for the adapter after resolving the
//     candidate location list
//   - target: the actual file chosen for this run after conflict handling;
//     it is either the primary file or the fallback file
func generateAdapterFile(adapterName string, projectRoot string, template []byte, candidates []string, opt Options) (Result, error) {
	rootAbs, err := absProjectRoot(projectRoot)
	if err != nil {
		return Result{}, err
	}
	primaryRel, err := resolvePrimaryLocation(adapterName, rootAbs, candidates)
	if err != nil {
		return Result{}, err
	}
	fallbackRel := fallbackPathFor(primaryRel)
	// Render once for the selected primary location. This is safe because the
	// fallback file lives beside the primary file, so both share the same
	// relative path to .context/ai_protocol.md.
	content := renderAdapterTemplate(template, primaryRel)

	out := writerOrStdout(opt.Writer)
	primaryPath := filepath.Join(rootAbs, primaryRel)
	fallbackPath := filepath.Join(rootAbs, fallbackRel)
	targetPath := primaryPath
	usedFallback := false
	message := ""

	// An existing primary file is always preserved. When it already exists, a
	// sidecar file is generated instead so the user can merge it manually.
	primaryExists, err := fileExists(primaryPath)
	if err != nil {
		return Result{}, fmt.Errorf("check %s: %w", primaryPath, err)
	}
	if primaryExists {
		targetPath = fallbackPath
		usedFallback = true
		message = fmt.Sprintf(
			"existing %s detected; generated %s instead. Append or merge its contents into %s manually.",
			primaryPath,
			fallbackPath,
			primaryPath,
		)
	}

	// The chosen target may itself already exist, for example when rerunning the
	// adapter after a previous fallback file was generated.
	targetExists, err := fileExists(targetPath)
	if err != nil {
		return Result{}, fmt.Errorf("check %s: %w", targetPath, err)
	}

	// Dry-run reports the file that would be written without touching disk.
	if opt.DryRun {
		action := ActionDryRunGenerate
		if targetExists && !opt.Force {
			action = ActionDryRunSkip
		}
		fmt.Fprintf(out, "  [%s] %s\n", action, targetPath)
		if usedFallback {
			printAdapterNote(out, action, primaryPath, targetPath)
		}
		return Result{
			Action:        action,
			GeneratedPath: targetPath,
			Skipped:       action == ActionDryRunSkip,
			UsedFallback:  usedFallback,
			Message:       message,
		}, nil
	}

	// Without force, never overwrite an existing target file. For adapters this
	// means force can replace a previously generated fallback file, but it still
	// does not replace the user's primary tool entrypoint.
	if targetExists && !opt.Force {
		fmt.Fprintf(out, "  [%s] %s\n", ActionSkipped, targetPath)
		if usedFallback {
			printAdapterNote(out, ActionSkipped, primaryPath, targetPath)
		}
		return Result{
			Action:        ActionSkipped,
			GeneratedPath: targetPath,
			Skipped:       true,
			UsedFallback:  usedFallback,
			Message:       message,
		}, nil
	}

	// At this point the target is either new or explicitly allowed to be
	// overwritten.
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return Result{}, fmt.Errorf("create adapter directory %s: %w", filepath.Dir(targetPath), err)
	}
	if err := os.WriteFile(targetPath, content, 0o644); err != nil {
		return Result{}, fmt.Errorf("write adapter %s: %w", targetPath, err)
	}

	fmt.Fprintf(out, "  [%s] %s\n", ActionGenerated, targetPath)
	if usedFallback {
		printAdapterNote(out, ActionGenerated, primaryPath, targetPath)
	}

	return Result{
		Action:        ActionGenerated,
		GeneratedPath: targetPath,
		UsedFallback:  usedFallback,
		Message:       message,
	}, nil
}
