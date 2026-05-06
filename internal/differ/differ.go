// Package differ provides line-level diff utilities for comparing
// individual env values and producing human-readable change descriptions.
package differ

import (
	"fmt"
	"strings"
)

// ValueDiff represents a single value change between two environments.
type ValueDiff struct {
	Key   string
	Left  string
	Right string
}

// Describe returns a short human-readable description of the diff.
func (d ValueDiff) Describe() string {
	return fmt.Sprintf("%s: %q → %q", d.Key, d.Left, d.Right)
}

// IsEmpty reports whether both sides are identical.
func (d ValueDiff) IsEmpty() bool {
	return d.Left == d.Right
}

// DiffValues compares two string maps and returns a slice of ValueDiff
// entries for keys whose values differ. Only keys present in both maps
// are considered; missing keys are ignored (see comparator for that).
func DiffValues(left, right map[string]string) []ValueDiff {
	var diffs []ValueDiff
	for key, lv := range left {
		if rv, ok := right[key]; ok && lv != rv {
			diffs = append(diffs, ValueDiff{
				Key:   key,
				Left:  lv,
				Right: rv,
			})
		}
	}
	sortDiffs(diffs)
	return diffs
}

// Summary returns a compact multi-line string listing all diffs.
func Summary(diffs []ValueDiff) string {
	if len(diffs) == 0 {
		return "no value differences"
	}
	lines := make([]string, len(diffs))
	for i, d := range diffs {
		lines[i] = d.Describe()
	}
	return strings.Join(lines, "\n")
}

// sortDiffs sorts diffs by key for deterministic output.
func sortDiffs(diffs []ValueDiff) {
	for i := 1; i < len(diffs); i++ {
		for j := i; j > 0 && diffs[j].Key < diffs[j-1].Key; j-- {
			diffs[j], diffs[j-1] = diffs[j-1], diffs[j]
		}
	}
}
