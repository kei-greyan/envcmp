package envgroup

import (
	"fmt"
	"sort"
	"strings"
)

// Format returns a human-readable text representation of the grouped result.
func Format(r Result) string {
	if len(r.Groups) == 0 {
		return "(no keys)"
	}
	var sb strings.Builder
	for _, label := range r.Order {
		bucket := r.Groups[label]
		fmt.Fprintf(&sb, "[%s] (%d keys)\n", label, len(bucket))
		keys := sortedKeys(bucket)
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s=%s\n", k, bucket[k])
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Summary returns a one-line summary of the grouping.
func Summary(r Result) string {
	total := 0
	for _, bucket := range r.Groups {
		total += len(bucket)
	}
	return fmt.Sprintf("%d group(s), %d key(s) total", len(r.Groups), total)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
