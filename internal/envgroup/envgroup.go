// Package envgroup groups environment variables by a shared delimiter pattern,
// producing named buckets suitable for structured inspection or export.
package envgroup

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how grouping is performed.
type Options struct {
	// Delimiter separates the group prefix from the rest of the key (default "_").
	Delimiter string
	// MaxDepth is the number of delimiter-separated segments used as the group key.
	// 0 or 1 means only the first segment is used.
	MaxDepth int
	// UngroupedLabel is the bucket name for keys with no delimiter (default "other").
	UngroupedLabel string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Delimiter:      "_",
		MaxDepth:       1,
		UngroupedLabel: "other",
	}
}

// Result holds the grouped output.
type Result struct {
	// Groups maps group label -> key -> value.
	Groups map[string]map[string]string
	// Order preserves the sorted group names for deterministic output.
	Order []string
}

// Group partitions env into named buckets according to opts.
func Group(env map[string]string, opts Options) (Result, error) {
	if opts.Delimiter == "" {
		return Result{}, fmt.Errorf("envgroup: delimiter must not be empty")
	}
	if opts.UngroupedLabel == "" {
		opts.UngroupedLabel = "other"
	}
	depth := opts.MaxDepth
	if depth < 1 {
		depth = 1
	}

	groups := make(map[string]map[string]string)

	for k, v := range env {
		label := extractLabel(k, opts.Delimiter, depth, opts.UngroupedLabel)
		if groups[label] == nil {
			groups[label] = make(map[string]string)
		}
		groups[label][k] = v
	}

	order := make([]string, 0, len(groups))
	for g := range groups {
		order = append(order, g)
	}
	sort.Strings(order)

	return Result{Groups: groups, Order: order}, nil
}

func extractLabel(key, delimiter string, depth int, ungrouped string) string {
	parts := strings.SplitN(key, delimiter, depth+1)
	if len(parts) <= 1 {
		return ungrouped
	}
	return strings.Join(parts[:depth], delimiter)
}
