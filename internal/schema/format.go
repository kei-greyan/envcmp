package schema

import (
	"fmt"
	"sort"
	"strings"
)

// Format returns a human-readable summary of schema violations.
func Format(violations []Violation) string {
	if len(violations) == 0 {
		return "schema: no violations found"
	}
	// Sort by key for deterministic output.
	sorted := make([]Violation, len(violations))
	copy(sorted, violations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("schema: %d violation(s) found\n", len(sorted)))
	for _, v := range sorted {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", v.Key, v.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Summary returns a single-line summary suitable for CLI output.
func Summary(violations []Violation) string {
	if len(violations) == 0 {
		return "OK"
	}
	return fmt.Sprintf("%d schema violation(s)", len(violations))
}
