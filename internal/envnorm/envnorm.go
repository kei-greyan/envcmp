// Package envnorm provides normalization utilities for environment variable maps.
// It handles trimming whitespace, normalizing key casing, and deduplicating entries.
package envnorm

import (
	"strings"
)

// Options controls normalization behaviour.
type Options struct {
	// TrimValues trims leading/trailing whitespace from values.
	TrimValues bool
	// LowercaseKeys converts all keys to lowercase.
	LowercaseKeys bool
	// UppercaseKeys converts all keys to uppercase. Takes precedence over LowercaseKeys.
	UppercaseKeys bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		TrimValues:    true,
		LowercaseKeys: false,
		UppercaseKeys: false,
	}
}

// Normalize applies the given Options to env and returns a new normalized map.
// The original map is never mutated.
func Normalize(env map[string]string, opts Options) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		normKey := normalizeKey(k, opts)
		normVal := normalizeValue(v, opts)
		result[normKey] = normVal
	}
	return result
}

// Diff returns keys whose normalized form differs from their original form.
// Useful for detecting keys that would be affected by normalization.
func Diff(env map[string]string, opts Options) map[string]string {
	changed := make(map[string]string)
	for k, v := range env {
		nk := normalizeKey(k, opts)
		nv := normalizeValue(v, opts)
		if nk != k || nv != v {
			changed[k] = nv
		}
	}
	return changed
}

func normalizeKey(k string, opts Options) string {
	if opts.UppercaseKeys {
		return strings.ToUpper(k)
	}
	if opts.LowercaseKeys {
		return strings.ToLower(k)
	}
	return k
}

func normalizeValue(v string, opts Options) string {
	if opts.TrimValues {
		return strings.TrimSpace(v)
	}
	return v
}
