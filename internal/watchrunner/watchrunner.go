// Package watchrunner wires the watcher to the full comparison pipeline so
// that results are re-printed to stdout whenever an .env file changes.
package watchrunner

import (
	"fmt"
	"io"
	"time"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/formatter"
	"github.com/user/envcmp/internal/parser"
	"github.com/user/envcmp/internal/watcher"
)

// Config holds the settings needed to run the watch loop.
type Config struct {
	LeftFile     string
	RightFile    string
	Format       string
	PollInterval time.Duration
	Out          io.Writer
}

// Run starts the watch loop. It performs an initial comparison immediately
// and then re-runs it each time either file changes. It blocks until done
// is closed.
func Run(cfg Config, done <-chan struct{}) error {
	if cfg.Out == nil {
		return fmt.Errorf("watchrunner: output writer must not be nil")
	}

	runComparison := func() {
		left, err := parser.ParseFile(cfg.LeftFile)
		if err != nil {
			fmt.Fprintf(cfg.Out, "[envcmp] error reading %s: %v\n", cfg.LeftFile, err)
			return
		}
		right, err := parser.ParseFile(cfg.RightFile)
		if err != nil {
			fmt.Fprintf(cfg.Out, "[envcmp] error reading %s: %v\n", cfg.RightFile, err)
			return
		}

		result := comparator.Compare(left, right)

		output, err := formatter.Render(result, cfg.Format)
		if err != nil {
			fmt.Fprintf(cfg.Out, "[envcmp] render error: %v\n", err)
			return
		}

		fmt.Fprintln(cfg.Out, "[envcmp] --- comparison updated ---")
		fmt.Fprint(cfg.Out, output)
	}

	// Run once immediately.
	runComparison()

	opts := watcher.Options{PollInterval: cfg.PollInterval}
	if opts.PollInterval <= 0 {
		opts = watcher.DefaultOptions()
	}

	return watcher.Watch(cfg.LeftFile, cfg.RightFile, opts, runComparison, done)
}
