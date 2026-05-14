// Package envsplit provides utilities for splitting a flat env map
// into named groups based on key prefix conventions (e.g. DB_, APP_, AWS_).
package envsplit

import (
	"fmt"
	"sort"
	"strings"
)

// Group represents a named subset of env keys sharing a common prefix.
type Group struct {
	Prefix string
	Keys   map[string]string
}

// Result holds all groups produced by a Split operation.
type Result struct {
	Groups   []Group
	Ungrouped map[string]string
}

// Split partitions env into groups by the provided prefixes.
// Keys not matching any prefix are placed in Ungrouped.
// Prefixes are matched case-sensitively and longest match wins.
func Split(env map[string]string, prefixes []string) Result {
	// Sort prefixes longest-first so longest match wins.
	sorted := make([]string, len(prefixes))
	copy(sorted, prefixes)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) > len(sorted[j])
	})

	groupMap := make(map[string]map[string]string, len(sorted))
	for _, p := range sorted {
		groupMap[p] = make(map[string]string)
	}
	ungrouped := make(map[string]string)

	for k, v := range env {
		matched := false
		for _, p := range sorted {
			if strings.HasPrefix(k, p) {
				groupMap[p][k] = v
				matched = true
				break
			}
		}
		if !matched {
			ungrouped[k] = v
		}
	}

	groups := make([]Group, 0, len(sorted))
	for _, p := range sorted {
		groups = append(groups, Group{Prefix: p, Keys: groupMap[p]})
	}
	// Return groups in original prefix order (re-sort by prefix string for stability).
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})

	return Result{Groups: groups, Ungrouped: ungrouped}
}

// Format returns a human-readable summary of the split result.
func Format(r Result) string {
	var sb strings.Builder
	for _, g := range r.Groups {
		keys := sortedKeys(g.Keys)
		sb.WriteString(fmt.Sprintf("[%s] %d key(s)\n", g.Prefix, len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", k, g.Keys[k]))
		}
	}
	if len(r.Ungrouped) > 0 {
		keys := sortedKeys(r.Ungrouped)
		sb.WriteString(fmt.Sprintf("[ungrouped] %d key(s)\n", len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", k, r.Ungrouped[k]))
		}
	}
	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
