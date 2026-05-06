// Package merger provides functionality to merge two parsed .env maps,
// producing a unified map with optional conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy int

const (
	// PreferLeft keeps the value from the left map when a key exists in both.
	PreferLeft Strategy = iota
	// PreferRight keeps the value from the right map when a key exists in both.
	PreferRight
	// ErrorOnConflict returns an error if the same key has different values.
	ErrorOnConflict
)

// Result holds the merged environment map and metadata about the merge.
type Result struct {
	Merged    map[string]string
	Conflicts []string // keys that had differing values
}

// Merge combines left and right env maps using the given strategy.
// Keys present in only one map are always included in the output.
func Merge(left, right map[string]string, strategy Strategy) (Result, error) {
	merged := make(map[string]string, len(left))
	var conflicts []string

	for k, v := range left {
		merged[k] = v
	}

	for k, rv := range right {
		lv, exists := merged[k]
		if !exists {
			merged[k] = rv
			continue
		}
		if lv == rv {
			continue
		}
		// Conflict: same key, different value.
		conflicts = append(conflicts, k)
		switch strategy {
		case PreferLeft:
			// keep existing lv — no-op
		case PreferRight:
			merged[k] = rv
		case ErrorOnConflict:
			return Result{}, fmt.Errorf("merger: conflict on key %q: %q vs %q", k, lv, rv)
		}
	}

	return Result{Merged: merged, Conflicts: conflicts}, nil
}
