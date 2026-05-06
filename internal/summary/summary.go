// Package summary provides a human-readable summary line for a comparison result.
package summary

import (
	"fmt"

	"github.com/user/envcmp/internal/stats"
)

// Summary holds a formatted summary of a comparison.
type Summary struct {
	TotalKeys    int
	Matched      int
	MissingLeft  int
	MissingRight int
	Mismatched   int
	Clean        bool
}

// Build constructs a Summary from a Stats value.
func Build(s stats.Stats) Summary {
	return Summary{
		TotalKeys:    s.Total,
		Matched:      s.Matched,
		MissingLeft:  s.MissingInLeft,
		MissingRight: s.MissingInRight,
		Mismatched:   s.Mismatched,
		Clean:        s.MissingInLeft == 0 && s.MissingInRight == 0 && s.Mismatched == 0,
	}
}

// OneLiner returns a concise single-line summary string.
func (s Summary) OneLiner() string {
	if s.Clean {
		return fmt.Sprintf("✓ %d keys matched, no differences found.", s.Matched)
	}
	return fmt.Sprintf(
		"✗ %d total keys: %d matched, %d missing in left, %d missing in right, %d mismatched.",
		s.TotalKeys,
		s.Matched,
		s.MissingLeft,
		s.MissingRight,
		s.Mismatched,
	)
}
