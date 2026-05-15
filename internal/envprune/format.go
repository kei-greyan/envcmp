package envprune

import (
	"fmt"
	"sort"
	"strings"
)

// Format returns a human-readable summary of a prune Result.
func Format(r Result) string {
	var sb strings.Builder

	prunedKeys := sortedKeys(r.Pruned)
	keptKeys := sortedKeys(r.Kept)

	if len(prunedKeys) == 0 {
		sb.WriteString("no keys pruned\n")
	} else {
		sb.WriteString("pruned keys:\n")
		for _, k := range prunedKeys {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	sb.WriteString(fmt.Sprintf("kept: %d  pruned: %d\n", len(keptKeys), len(prunedKeys)))
	return sb.String()
}

// Summary returns a one-line summary string.
func Summary(r Result) string {
	return fmt.Sprintf("%d kept, %d pruned", len(r.Kept), len(r.Pruned))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
