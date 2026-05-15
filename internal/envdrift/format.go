package envdrift

import (
	"fmt"
	"strings"
)

// Format renders a Result as a human-readable text report.
func Format(r Result) string {
	if !r.HasDrift() {
		return "No drift detected.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n", Summary(r)))
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	for _, e := range r.Entries {
		switch e.Kind {
		case DriftAdded:
			sb.WriteString(fmt.Sprintf("+ %-30s (added)\n", e.Key))
		case DriftRemoved:
			sb.WriteString(fmt.Sprintf("- %-30s (removed)\n", e.Key))
		case DriftChanged:
			sb.WriteString(fmt.Sprintf("~ %-30s baseline=%q current=%q\n", e.Key, e.Baseline, e.Current))
		}
	}

	return sb.String()
}
