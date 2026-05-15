// Package envscope provides scoped views of an env map by filtering keys
// that match a given prefix or suffix pattern, returning a namespaced subset.
package envscope

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how the scope is extracted.
type Options struct {
	// Prefix filters keys that start with this string.
	Prefix string
	// Suffix filters keys that end with this string.
	Suffix string
	// StripPrefix removes the prefix from the resulting keys when true.
	StripPrefix bool
	// StripSuffix removes the suffix from the resulting keys when true.
	StripSuffix bool
}

// Result holds the scoped env map and metadata.
type Result struct {
	Env      map[string]string
	Matched  int
	Excluded int
}

// Extract returns a scoped view of env according to opts.
// At least one of Prefix or Suffix must be set.
func Extract(env map[string]string, opts Options) (Result, error) {
	if opts.Prefix == "" && opts.Suffix == "" {
		return Result{}, fmt.Errorf("envscope: at least one of Prefix or Suffix must be set")
	}

	scoped := make(map[string]string)
	excluded := 0

	for k, v := range env {
		if !matches(k, opts) {
			excluded++
			continue
		}
		key := transform(k, opts)
		scoped[key] = v
	}

	return Result{
		Env:      scoped,
		Matched:  len(scoped),
		Excluded: excluded,
	}, nil
}

// Keys returns the sorted keys of a scoped result.
func Keys(r Result) []string {
	keys := make([]string, 0, len(r.Env))
	for k := range r.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func matches(key string, opts Options) bool {
	if opts.Prefix != "" && !strings.HasPrefix(key, opts.Prefix) {
		return false
	}
	if opts.Suffix != "" && !strings.HasSuffix(key, opts.Suffix) {
		return false
	}
	return true
}

func transform(key string, opts Options) string {
	if opts.StripPrefix && opts.Prefix != "" {
		key = strings.TrimPrefix(key, opts.Prefix)
	}
	if opts.StripSuffix && opts.Suffix != "" {
		key = strings.TrimSuffix(key, opts.Suffix)
	}
	return key
}
