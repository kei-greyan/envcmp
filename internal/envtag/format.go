package envtag

import (
	"fmt"
	"strings"
)

// Format renders a Result as a human-readable text block.
func Format(r Result) string {
	var sb strings.Builder

	for _, tag := range TagNames(r) {
		keys := r.Tagged[tag]
		fmt.Fprintf(&sb, "[%s] (%d keys)\n", tag, len(keys))
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s\n", k)
		}
	}

	if len(r.Untagged) > 0 {
		fmt.Fprintf(&sb, "[untagged] (%d keys)\n", len(r.Untagged))
		for _, k := range r.Untagged {
			fmt.Fprintf(&sb, "  %s\n", k)
		}
	}

	return sb.String()
}

// Summary returns a one-line summary of the tagging result.
func Summary(r Result) string {
	total := len(r.Untagged)
	for _, keys := range r.Tagged {
		total += len(keys)
	}
	return fmt.Sprintf("%d tag(s), %d tagged key(s), %d untagged key(s)",
		len(r.Tagged), total-len(r.Untagged), len(r.Untagged))
}
