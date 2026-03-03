package render

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/thyrlian/ctx-init/internal/plan"
)

// Action represents the outcome of processing a single entry
type Action int

const (
	ActionUnknown    Action = iota // zero value; returned alongside a non-nil error
	ActionCopied                   // file was copied to destination
	ActionSkipped                  // destination already exists, not overwritten
	ActionDryRunCopy               // would copy (dry-run)
	ActionDryRunSkip               // would skip (dry-run)
)

func (a Action) String() string {
	switch a {
	case ActionCopied:
		return "copied"
	case ActionSkipped:
		return "skipped"
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
	Copied  int
	Skipped int
}

func (r *Result) record(a Action) {
	switch a {
	case ActionCopied, ActionDryRunCopy:
		r.Copied++
	case ActionSkipped, ActionDryRunSkip:
		r.Skipped++
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
		action, err := processEntry(e, opt)
		if err != nil {
			return res, err
		}
		fmt.Fprintf(w, "  [%s] %s\n", action, e.Dst)
		res.record(action)
	}
	return res, nil
}

func processEntry(e plan.Entry, opt Options) (Action, error) {
	exists, err := fileExists(e.Dst)
	if err != nil {
		return ActionUnknown, fmt.Errorf("check destination %s: %w", e.Dst, err)
	}

	if opt.DryRun {
		if exists && !opt.Force {
			return ActionDryRunSkip, nil
		}
		return ActionDryRunCopy, nil
	}

	if exists && !opt.Force {
		return ActionSkipped, nil
	}

	dir := filepath.Dir(e.Dst)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ActionUnknown, fmt.Errorf("mkdir %s: %w", dir, err)
	}

	if err := copyFile(e.Src, e.Dst); err != nil {
		return ActionUnknown, fmt.Errorf("copy %s → %s: %w", e.Src, e.Dst, err)
	}

	return ActionCopied, nil
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
