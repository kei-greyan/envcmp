package envfreeze

import (
	"fmt"
	"strings"
)

// Format renders the list of violations as a human-readable text report.
func Format(violations []Violation) string {
	if len(violations) == 0 {
		return "no mutations detected against frozen snapshot\n"
	}

	var sb strings.Builder
	sb.WriteString("frozen snapshot violations:\n")

	for _, v := range violations {
		switch v.Kind {
		case "modified":
			sb.WriteString(fmt.Sprintf("  [modified] %s: %q -> %q\n", v.Key, v.Frozen, v.Current))
		case "added":
			sb.WriteString(fmt.Sprintf("  [added]    %s: (absent) -> %q\n", v.Key, v.Current))
		case "removed":
			sb.WriteString(fmt.Sprintf("  [removed]  %s: %q -> (absent)\n", v.Key, v.Frozen))
		}
	}

	sb.WriteString(fmt.Sprintf("\n%s\n", Summary(violations)))
	return sb.String()
}
