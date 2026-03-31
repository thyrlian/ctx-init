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
)

type generator func(projectRoot string, opt Options) (Result, error)

var supportedAdapters = []string{
	AdapterClaude,
}

var generators = map[string]generator{
	AdapterClaude: generateClaude,
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

// generateAdapterFile applies the standard ctx-init adapter behavior for a
// single output file:
//   - write the primary file when it does not exist
//   - fall back to an alternate filename when the primary file already exists
//   - respect dry-run semantics and only use force to overwrite an existing
//     fallback file; an existing primary file is never overwritten implicitly
//   - print a user-facing note when manual merge is required
//
// This keeps the common "one adapter file plus fallback filename" flow in one
// place so individual adapters only need to supply their content and names.
// The provided names must be project-root-relative file paths.
//
// Terminology:
//   - primary: the preferred output filename for the adapter
//   - target: the actual file chosen for this run after conflict handling;
//     it is either the primary file or the fallback file
func generateAdapterFile(projectRoot string, content []byte, primaryName string, fallbackName string, opt Options) (Result, error) {
	rootAbs, err := absProjectRoot(projectRoot)
	if err != nil {
		return Result{}, err
	}
	primaryRel, err := projectRelativePath(primaryName)
	if err != nil {
		return Result{}, fmt.Errorf("invalid primary path: %w", err)
	}
	fallbackRel, err := projectRelativePath(fallbackName)
	if err != nil {
		return Result{}, fmt.Errorf("invalid fallback path: %w", err)
	}

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
		fmt.Fprintf(out, "  [dry-run/generate] %s\n", targetPath)
		if message != "" {
			fmt.Fprintf(out, "  note: %s\n", message)
		}
		return Result{
			GeneratedPath: targetPath,
			UsedFallback:  usedFallback,
			Message:       message,
		}, nil
	}

	// Without force, never overwrite an existing target file. For adapters this
	// means force can replace a previously generated fallback file, but it still
	// does not replace the user's primary tool entrypoint.
	if targetExists && !opt.Force {
		fmt.Fprintf(out, "  [skipped] %s\n", targetPath)
		if message != "" {
			fmt.Fprintf(out, "  note: %s\n", message)
		}
		return Result{
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

	fmt.Fprintf(out, "  [generated] %s\n", targetPath)
	if message != "" {
		fmt.Fprintf(out, "  note: %s\n", message)
	}

	return Result{
		GeneratedPath: targetPath,
		UsedFallback:  usedFallback,
		Message:       message,
	}, nil
}
