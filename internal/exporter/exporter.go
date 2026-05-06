// Package exporter writes comparison results to files or stdout
// in a specified format, enabling piping and file-based reporting.
package exporter

import (
	"fmt"
	"os"

	"github.com/user/envcmp/internal/formatter"
	"github.com/user/envcmp/internal/comparator"
)

// Options controls how the result is exported.
type Options struct {
	// OutputFile is the path to write output to. If empty, writes to stdout.
	OutputFile string
	// Format is the output format: text, json, or markdown.
	Format string
	// Overwrite controls whether an existing file is overwritten.
	Overwrite bool
}

// Export renders the comparison result and writes it to the configured destination.
func Export(result comparator.Result, opts Options) error {
	output, err := formatter.Render(result, opts.Format)
	if err != nil {
		return fmt.Errorf("exporter: render failed: %w", err)
	}

	if opts.OutputFile == "" {
		_, err = fmt.Fprint(os.Stdout, output)
		return err
	}

	return writeFile(opts.OutputFile, output, opts.Overwrite)
}

// writeFile writes content to a file, respecting the overwrite flag.
func writeFile(path, content string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("exporter: file already exists: %s (use --overwrite to replace)", path)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("exporter: could not create file %s: %w", path, err)
	}
	defer f.Close()

	_, err = fmt.Fprint(f, content)
	if err != nil {
		return fmt.Errorf("exporter: write failed: %w", err)
	}
	return nil
}
