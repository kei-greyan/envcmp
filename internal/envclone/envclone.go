// Package envclone provides utilities for cloning an env map with optional
// key renaming, prefix stripping, and prefix injection.
package envclone

import (
	"fmt"
	"strings"
)

// Options controls how the clone operation behaves.
type Options struct {
	// StripPrefix removes this prefix from every key before writing to the clone.
	StripPrefix string
	// AddPrefix prepends this prefix to every key in the clone.
	AddPrefix string
	// Only includes keys whose original name appears in this set.
	// When nil or empty all keys are included.
	OnlyKeys map[string]struct{}
}

// Result holds the cloned environment and a log of every transformation.
type Result struct {
	Env     map[string]string
	Entries []Entry
}

// Entry records a single key transformation.
type Entry struct {
	OriginalKey string
	ClonedKey   string
	Value       string
	Skipped     bool
}

// Clone copies src into a new map, applying the transformations in opts.
// It returns an error if a key collision occurs after renaming.
func Clone(src map[string]string, opts Options) (Result, error) {
	out := make(map[string]string, len(src))
	entries := make([]Entry, 0, len(src))

	for k, v := range src {
		entry := Entry{OriginalKey: k, Value: v}

		if len(opts.OnlyKeys) > 0 {
			if _, ok := opts.OnlyKeys[k]; !ok {
				entry.Skipped = true
				entries = append(entries, entry)
				continue
			}
		}

		newKey := k
		if opts.StripPrefix != "" {
			newKey = strings.TrimPrefix(newKey, opts.StripPrefix)
		}
		if opts.AddPrefix != "" {
			newKey = opts.AddPrefix + newKey
		}

		if _, exists := out[newKey]; exists {
			return Result{}, fmt.Errorf("envclone: key collision on %q (from %q)", newKey, k)
		}

		out[newKey] = v
		entry.ClonedKey = newKey
		entries = append(entries, entry)
	}

	return Result{Env: out, Entries: entries}, nil
}
