// Package sorter provides utilities for sorting comparison results
// for deterministic output across runs.
package sorter

import (
	"sort"

	"github.com/yourusername/envcmp/internal/comparator"
)

// SortedResult holds sorted slices derived from a comparator.Result
// so that output is stable regardless of map iteration order.
type SortedResult struct {
	MissingInRight []string
	MissingInLeft  []string
	Mismatched     []MismatchEntry
}

// MismatchEntry holds a key and its two differing values.
type MismatchEntry struct {
	Key    string
	Left   string
	Right  string
}

// Sort takes a comparator.Result and returns a SortedResult with all
// keys and mismatch entries sorted alphabetically.
func Sort(r comparator.Result) SortedResult {
	missingRight := make([]string, 0, len(r.MissingInRight))
	for k := range r.MissingInRight {
		missingRight = append(missingRight, k)
	}
	sort.Strings(missingRight)

	missingLeft := make([]string, 0, len(r.MissingInLeft))
	for k := range r.MissingInLeft {
		missingLeft = append(missingLeft, k)
	}
	sort.Strings(missingLeft)

	mismatched := make([]MismatchEntry, 0, len(r.Mismatched))
	for k, v := range r.Mismatched {
		mismatched = append(mismatched, MismatchEntry{
			Key:   k,
			Left:  v.Left,
			Right: v.Right,
		})
	}
	sort.Slice(mismatched, func(i, j int) bool {
		return mismatched[i].Key < mismatched[j].Key
	})

	return SortedResult{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Mismatched:     mismatched,
	}
}
