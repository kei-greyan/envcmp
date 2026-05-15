// Package envflatten provides utilities for flattening nested key structures
// in environment maps using a configurable delimiter.
package envflatten

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how flattening is performed.
type Options struct {
	// Delimiter separates nested key segments. Defaults to "__".
	Delimiter string
	// Uppercase converts all output keys to uppercase.
	Uppercase bool
}

// DefaultOptions returns sensible defaults for flattening.
func DefaultOptions() Options {
	return Options{
		Delimiter: "__",
		Uppercase: false,
	}
}

// FlatEntry represents a single flattened key-value pair with its origin prefix.
type FlatEntry struct {
	Key    string
	Value  string
	Prefix string
}

// Flatten takes an env map and groups keys by their prefix (first segment split
// by the delimiter), returning a map of prefix -> slice of FlatEntry.
// Keys without a delimiter are placed under the empty-string prefix.
func Flatten(env map[string]string, opts Options) map[string][]FlatEntry {
	if opts.Delimiter == "" {
		opts.Delimiter = "__"
	}

	result := make(map[string][]FlatEntry)

	for k, v := range env {
		key := k
		if opts.Uppercase {
			key = strings.ToUpper(key)
		}

		parts := strings.SplitN(key, opts.Delimiter, 2)
		prefix := ""
		if len(parts) == 2 {
			prefix = parts[0]
		}

		result[prefix] = append(result[prefix], FlatEntry{
			Key:    key,
			Value:  v,
			Prefix: prefix,
		})
	}

	for prefix := range result {
		sort.Slice(result[prefix], func(i, j int) bool {
			return result[prefix][i].Key < result[prefix][j].Key
		})
	}

	return result
}

// Format returns a human-readable text representation of flattened groups.
func Format(groups map[string][]FlatEntry) string {
	if len(groups) == 0 {
		return "(empty)\n"
	}

	prefixes := make([]string, 0, len(groups))
	for p := range groups {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	var sb strings.Builder
	for _, p := range prefixes {
		header := p
		if header == "" {
			header = "(no prefix)"
		}
		sb.WriteString(fmt.Sprintf("[%s]\n", header))
		for _, e := range groups[p] {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, e.Value))
		}
	}
	return sb.String()
}

// Summary returns a one-line summary of the flattened result.
func Summary(groups map[string][]FlatEntry) string {
	total := 0
	for _, entries := range groups {
		total += len(entries)
	}
	return fmt.Sprintf("%d key(s) across %d group(s)", total, len(groups))
}
