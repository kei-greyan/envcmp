package envrotate

import (
	"fmt"
	"strings"
)

// Format returns a human-readable summary of the rotation result.
func Format(r Result) string {
	if len(r.Rotated) == 0 {
		return "No keys were rotated."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Rotated %d key(s):\n", len(r.Rotated)))
	for _, e := range r.Rotated {
		oldDisplay := maskValue(e.OldValue)
		newDisplay := maskValue(e.NewValue)
		sb.WriteString(fmt.Sprintf("  %-30s %s -> %s\n", e.Key, oldDisplay, newDisplay))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Summary returns a one-line summary suitable for logging.
func Summary(r Result) string {
	if len(r.Rotated) == 0 {
		return "rotation complete: 0 keys rotated"
	}
	keys := make([]string, 0, len(r.Rotated))
	for _, e := range r.Rotated {
		keys = append(keys, e.Key)
	}
	return fmt.Sprintf("rotation complete: %d key(s) rotated (%s)", len(r.Rotated), strings.Join(keys, ", "))
}

func maskValue(v string) string {
	if v == "" {
		return "(empty)"
	}
	if len(v) <= 4 {
		return "****"
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
