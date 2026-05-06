// Package stats provides summary statistics derived from a comparison result.
package stats

import "github.com/user/envcmp/internal/comparator"

// Summary holds aggregated counts from a comparison result.
type Summary struct {
	TotalKeys      int
	MissingInLeft  int
	MissingInRight int
	Mismatched     int
	Matched        int
}

// Compute derives a Summary from the given comparator.Result.
// TotalKeys is the union of all unique keys observed across both sides.
func Compute(r comparator.Result) Summary {
	seen := make(map[string]struct{})

	for _, k := range r.MissingInLeft {
		seen[k] = struct{}{}
	}
	for _, k := range r.MissingInRight {
		seen[k] = struct{}{}
	}
	for k := range r.Mismatched {
		seen[k] = struct{}{}
	}
	for k := range r.Matched {
		seen[k] = struct{}{}
	}

	s := Summary{
		TotalKeys:      len(seen),
		MissingInLeft:  len(r.MissingInLeft),
		MissingInRight: len(r.MissingInRight),
		Mismatched:     len(r.Mismatched),
		Matched:        len(r.Matched),
	}
	return s
}

// HasDiff returns true when the summary contains any discrepancy.
func (s Summary) HasDiff() bool {
	return s.MissingInLeft > 0 || s.MissingInRight > 0 || s.Mismatched > 0
}
