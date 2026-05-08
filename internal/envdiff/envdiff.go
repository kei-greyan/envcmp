// Package envdiff provides line-level diffing between two versions of the same
// .env file, showing added, removed, and changed lines.
package envdiff

import (
	"fmt"
	"sort"
	"strings"
)

// LineKind describes the type of change for a diff entry.
type LineKind string

const (
	Added   LineKind = "added"
	Removed LineKind = "removed"
	Changed LineKind = "changed"
	Same    LineKind = "same"
)

// Entry represents a single key-level diff entry between two env maps.
type Entry struct {
	Key      string
	Kind     LineKind
	OldValue string
	NewValue string
}

// Result holds all diff entries produced by Diff.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if any entry is not Same.
func (r Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Kind != Same {
			return true
		}
	}
	return false
}

// Diff computes a line-level diff between two env maps (old vs new).
// All keys from both maps are included in the result, sorted alphabetically.
func Diff(oldEnv, newEnv map[string]string) Result {
	keys := unionKeys(oldEnv, newEnv)
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		oldVal, inOld := oldEnv[k]
		newVal, inNew := newEnv[k]

		var kind LineKind
		switch {
		case inOld && !inNew:
			kind = Removed
		case !inOld && inNew:
			kind = Added
		case oldVal != newVal:
			kind = Changed
		default:
			kind = Same
		}

		entries = append(entries, Entry{
			Key:      k,
			Kind:     kind,
			OldValue: oldVal,
			NewValue: newVal,
		})
	}
	return Result{Entries: entries}
}

// Format renders the diff result as a human-readable string.
func Format(r Result) string {
	if !r.HasChanges() {
		return "no changes detected"
	}
	var sb strings.Builder
	for _, e := range r.Entries {
		switch e.Kind {
		case Added:
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, e.NewValue))
		case Removed:
			sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, e.OldValue))
		case Changed:
			sb.WriteString(fmt.Sprintf("~ %s: %s → %s\n", e.Key, e.OldValue, e.NewValue))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
