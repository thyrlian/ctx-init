package render

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/thyrlian/ctx-init/internal/plan"
)

// ctxIDPattern matches a ctx-id comment line (multiline mode, anchored to line boundaries).
var ctxIDPattern = regexp.MustCompile(`(?m)^\s*<!-- ctx-id: [0-9a-f]{16} -->\s*$`)

// Action represents the outcome of processing a single entry
type Action int

const (
	ActionUnknown        Action = iota // zero value; returned alongside a non-nil error
	ActionGenerated                    // file was generated (e.g. _INDEX.md)
	ActionCopied                       // file was copied to destination
	ActionSkipped                      // destination already exists, not overwritten
	ActionDryRunGenerate               // would generate (dry-run)
	ActionDryRunCopy                   // would copy (dry-run)
	ActionDryRunSkip                   // would skip (dry-run)
)

func (a Action) String() string {
	switch a {
	case ActionGenerated:
		return "generated"
	case ActionCopied:
		return "copied"
	case ActionSkipped:
		return "skipped"
	case ActionDryRunGenerate:
		return "dry-run/generate"
	case ActionDryRunCopy:
		return "dry-run/copy"
	case ActionDryRunSkip:
		return "dry-run/skip"
	default:
		return "unknown"
	}
}

// Options controls render behaviour
type Options struct {
	DryRun bool
	Force  bool
	Writer io.Writer // output destination; defaults to os.Stdout if nil
}

// Result summarises what Run did (or would do in dry-run)
type Result struct {
	Generated int
	Copied    int
	Skipped   int
	FileIDs   map[string]string // rel slash path (relative to OutRootAbs) → ctx-id
}

func (r *Result) record(dst string, a Action, ctxID string) {
	switch a {
	case ActionGenerated, ActionDryRunGenerate:
		r.Generated++
	case ActionCopied, ActionDryRunCopy:
		r.Copied++
	case ActionSkipped, ActionDryRunSkip:
		r.Skipped++
	}
	if ctxID != "" {
		if r.FileIDs == nil {
			r.FileIDs = make(map[string]string)
		}
		r.FileIDs[dst] = ctxID
	}
}

// Run executes (or simulates) every entry in the plan
func Run(p *plan.Plan, opt Options) (Result, error) {
	w := opt.Writer
	if w == nil {
		w = os.Stdout
	}

	var res Result
	for _, e := range p.Entries {
		action, ctxID, err := processEntry(e, opt)
		if err != nil {
			return res, err
		}

		rel, err := filepath.Rel(p.OutRootAbs, e.Dst)
		if err != nil {
			return res, fmt.Errorf("rel path for %s: %w", e.Dst, err)
		}
		relSlash := filepath.ToSlash(rel)

		if ctxID != "" {
			fmt.Fprintf(w, "  [%s] %s   ctx-id: %s\n", action, e.Dst, ctxID)
		} else {
			fmt.Fprintf(w, "  [%s] %s\n", action, e.Dst)
		}
		res.record(relSlash, action, ctxID)
	}

	// Generate _INDEX.md (always, not part of manifest entries)
	indexAction, err := generateIndex(p, opt, w)
	if err != nil {
		return res, err
	}
	res.record("_INDEX.md", indexAction, "")

	return res, nil
}

func processEntry(e plan.Entry, opt Options) (Action, string, error) {
	exists, err := fileExists(e.Dst)
	if err != nil {
		return ActionUnknown, "", fmt.Errorf("check destination %s: %w", e.Dst, err)
	}

	if opt.DryRun {
		if exists && !opt.Force {
			return ActionDryRunSkip, "", nil
		}
		return ActionDryRunCopy, "", nil
	}

	if exists && !opt.Force {
		return ActionSkipped, "", nil
	}

	dir := filepath.Dir(e.Dst)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ActionUnknown, "", fmt.Errorf("mkdir %s: %w", dir, err)
	}

	if err := copyFile(e.Src, e.Dst); err != nil {
		return ActionUnknown, "", fmt.Errorf("copy %s → %s: %w", e.Src, e.Dst, err)
	}

	var ctxID string
	if strings.HasSuffix(e.Dst, ".md") {
		ctxID, err = appendCtxID(e.Dst)
		if err != nil {
			return ActionUnknown, "", fmt.Errorf("append ctx-id to %s: %w", e.Dst, err)
		}
	}

	return ActionCopied, ctxID, nil
}

// appendCtxID writes a random proof-of-read token to the file as a markdown comment.
// If a ctx-id already exists anywhere in the file, it is replaced rather than duplicated.
// The original file permissions are preserved.
func appendCtxID(dst string) (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	line := fmt.Sprintf("<!-- ctx-id: %s -->", token)

	info, err := os.Stat(dst)
	if err != nil {
		return "", err
	}
	mode := info.Mode()

	content, err := os.ReadFile(dst)
	if err != nil {
		return "", err
	}

	var newContent []byte
	if ctxIDPattern.Match(content) {
		// Replace existing ctx-id line in place; no surrounding newlines added
		newContent = ctxIDPattern.ReplaceAll(content, []byte(line))
	} else {
		// Append: ensure a blank line before the token for readability
		if len(content) > 0 && content[len(content)-1] != '\n' {
			newContent = append(content, '\n', '\n')
		} else {
			newContent = append(content, '\n')
		}
		newContent = append(newContent, []byte(line+"\n")...)
	}

	if err := os.WriteFile(dst, newContent, mode); err != nil {
		return "", err
	}
	return token, nil
}

// generateIndex generates _INDEX.md in the output root
func generateIndex(p *plan.Plan, opt Options, w io.Writer) (Action, error) {
	indexDst := filepath.Join(p.OutRootAbs, "_INDEX.md")

	if opt.DryRun {
		fmt.Fprintf(w, "  [dry-run/generate] %s\n", indexDst)
		return ActionDryRunGenerate, nil
	}

	if err := os.MkdirAll(p.OutRootAbs, 0o755); err != nil {
		return ActionUnknown, fmt.Errorf("mkdir %s: %w", p.OutRootAbs, err)
	}

	content, err := buildIndexContent(p)
	if err != nil {
		return ActionUnknown, err
	}

	if err := os.WriteFile(indexDst, []byte(content), 0o644); err != nil {
		return ActionUnknown, fmt.Errorf("write _INDEX.md: %w", err)
	}

	fmt.Fprintf(w, "  [generated] %s\n", indexDst)
	return ActionGenerated, nil
}

func buildIndexContent(p *plan.Plan) (string, error) {
	var sb strings.Builder

	sb.WriteString("# Context Index\n\n")
	sb.WriteString(fmt.Sprintf(
		"> Generated by ctx-init on %s | preset: `%s`\n",
		time.Now().Format("2006-01-02"), p.Preset,
	))
	sb.WriteString(">\n")
	sb.WriteString("> Tags after a file name indicate load priority — see `ai_protocol.md` for details.\n\n")
	sb.WriteString("## Files\n\n")

	for _, e := range p.Entries {
		rel, err := filepath.Rel(p.OutRootAbs, e.Dst)
		if err != nil {
			return "", fmt.Errorf("rel path for %s: %w", e.Dst, err)
		}
		relSlash := filepath.ToSlash(rel)

		var tagStrs []string
		for _, t := range e.Tags {
			tagStrs = append(tagStrs, "`"+t+"`")
		}

		if len(tagStrs) > 0 {
			sb.WriteString(fmt.Sprintf("- [%s](%s) %s\n", relSlash, relSlash, strings.Join(tagStrs, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("- [%s](%s)\n", relSlash, relSlash))
		}
	}

	return sb.String(), nil
}

// fileExists reports whether path exists
// Returns an error only for unexpected failures (e.g., permission denied)
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

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
