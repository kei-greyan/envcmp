package envmerge

import (
	"fmt"
	"strings"
)

// Format returns a human-readable summary of the merge result.
func Format(r Result) string {
	if len(r.Entries) == 0 {
		return "no keys merged\n"
	}

	var sb strings.Builder
	for _, e := range r.Entries {
		if e.Conflict {
			fmt.Fprintf(&sb, "~ %s (conflict, source %d wins)\n", e.Key, e.Source)
		} else {
			fmt.Fprintf(&sb, "  %s (source %d)\n", e.Key, e.Source)
		}
	}
	return sb.String()
}

// Summary returns a one-line description of the merge outcome.
func Summary(r Result) string {
	total := len(r.Entries)
	conflicts := 0
	for _, e := range r.Entries {
		if e.Conflict {
			conflicts++
		}
	}
	if conflicts == 0 {
		return fmt.Sprintf("%d keys merged, no conflicts", total)
	}
	return fmt.Sprintf("%d keys merged, %d conflict(s) resolved", total, conflicts)
}
