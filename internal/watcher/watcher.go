// Package watcher monitors .env files for changes and triggers a callback
// when a modification is detected.
package watcher

import (
	"fmt"
	"os"
	"time"
)

// Options configures the file watcher behaviour.
type Options struct {
	// PollInterval is how often the watcher checks for changes.
	PollInterval time.Duration
}

// DefaultOptions returns sensible defaults for the watcher.
func DefaultOptions() Options {
	return Options{
		PollInterval: 2 * time.Second,
	}
}

// Watch polls the two given file paths for changes and calls onChange each
// time either file is modified. It blocks until the done channel is closed.
func Watch(left, right string, opts Options, onChange func(), done <-chan struct{}) error {
	if left == "" || right == "" {
		return fmt.Errorf("watcher: both file paths must be non-empty")
	}

	if opts.PollInterval <= 0 {
		opts.PollInterval = DefaultOptions().PollInterval
	}

	lastLeft, err := modTime(left)
	if err != nil {
		return fmt.Errorf("watcher: cannot stat %s: %w", left, err)
	}
	lastRight, err := modTime(right)
	if err != nil {
		return fmt.Errorf("watcher: cannot stat %s: %w", right, err)
	}

	ticker := time.NewTicker(opts.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			newLeft, err := modTime(left)
			if err != nil {
				continue
			}
			newRight, err := modTime(right)
			if err != nil {
				continue
			}

			if newLeft != lastLeft || newRight != lastRight {
				lastLeft = newLeft
				lastRight = newRight
				onChange()
			}
		}
	}
}

func modTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}
