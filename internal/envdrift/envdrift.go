// Package envdrift detects configuration drift between a saved baseline
// and the current environment map.
package envdrift

import (
	"fmt"
	"sort"
)

// DriftKind describes the type of drift detected for a key.
type DriftKind string

const (
	DriftAdded   DriftKind = "added"   // present in current, absent in baseline
	DriftRemoved DriftKind = "removed" // present in baseline, absent in current
	DriftChanged DriftKind = "changed" // value differs between baseline and current
)

// Entry represents a single drifted key.
type Entry struct {
	Key      string
	Kind     DriftKind
	Baseline string // empty when Kind == DriftAdded
	Current  string // empty when Kind == DriftRemoved
}

// Result holds all drift entries detected between two environments.
type Result struct {
	Entries []Entry
}

// HasDrift returns true when at least one drift entry exists.
func (r Result) HasDrift() bool { return len(r.Entries) > 0 }

// Detect compares baseline and current env maps and returns all drift.
func Detect(baseline, current map[string]string) Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, bv := range baseline {
		seen[k] = true
		cv, ok := current[k]
		if !ok {
			entries = append(entries, Entry{Key: k, Kind: DriftRemoved, Baseline: bv})
		} else if bv != cv {
			entries = append(entries, Entry{Key: k, Kind: DriftChanged, Baseline: bv, Current: cv})
		}
	}

	for k, cv := range current {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Kind: DriftAdded, Current: cv})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Key != entries[j].Key {
			return entries[i].Key < entries[j].Key
		}
		return entries[i].Kind < entries[j].Kind
	})

	return Result{Entries: entries}
}

// Summary returns a one-line human-readable drift summary.
func Summary(r Result) string {
	if !r.HasDrift() {
		return "no drift detected"
	}
	var added, removed, changed int
	for _, e := range r.Entries {
		switch e.Kind {
		case DriftAdded:
			added++
		case DriftRemoved:
			removed++
		case DriftChanged:
			changed++
		}
	}
	return fmt.Sprintf("drift detected: %d added, %d removed, %d changed", added, removed, changed)
}
