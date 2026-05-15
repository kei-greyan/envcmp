package envhash

import (
	"fmt"
	"sort"
	"strings"
)

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("hash: %s\n(no keys)\n", r.Hash)
	}

	keys := make([]string, 0, len(r.Entries))
	for k := range r.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	fmt.Fprintf(&sb, "hash: %s\n", r.Hash)
	for _, k := range keys {
		fmt.Fprintf(&sb, "  %-30s %s\n", k, r.Entries[k][:12]+"...")
	}
	return sb.String()
}

// Summary returns a one-line description suitable for logging or CLI output.
func Summary(r Result) string {
	return fmt.Sprintf("envhash: %d key(s), hash=%s", len(r.Entries), r.Hash[:12]+"...")
}
