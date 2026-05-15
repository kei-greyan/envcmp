// Package envaudit provides audit trail functionality for tracking
// changes to environment variables over time.
package envaudit

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ChangeKind describes the type of change recorded in an audit entry.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
)

// Entry represents a single audited change to an environment key.
type Entry struct {
	Key       string     `json:"key"`
	Kind      ChangeKind `json:"kind"`
	OldValue  string     `json:"old_value,omitempty"`
	NewValue  string     `json:"new_value,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

// Report holds all audit entries produced by comparing two env snapshots.
type Report struct {
	Entries []Entry `json:"entries"`
}

// IsClean returns true when no changes were detected.
func (r Report) IsClean() bool {
	return len(r.Entries) == 0
}

// Audit compares two env maps (before, after) and returns an audit Report
// capturing every addition, removal, and modification.
func Audit(before, after map[string]string) Report {
	now := time.Now().UTC()
	var entries []Entry

	for key, oldVal := range before {
		newVal, exists := after[key]
		if !exists {
			entries = append(entries, Entry{Key: key, Kind: Removed, OldValue: oldVal, Timestamp: now})
		} else if newVal != oldVal {
			entries = append(entries, Entry{Key: key, Kind: Modified, OldValue: oldVal, NewValue: newVal, Timestamp: now})
		}
	}

	for key, newVal := range after {
		if _, exists := before[key]; !exists {
			entries = append(entries, Entry{Key: key, Kind: Added, NewValue: newVal, Timestamp: now})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Report{Entries: entries}
}

// Format returns a human-readable text representation of the audit report.
func Format(r Report) string {
	if r.IsClean() {
		return "audit: no changes detected\n"
	}
	var sb strings.Builder
	for _, e := range r.Entries {
		switch e.Kind {
		case Added:
			fmt.Fprintf(&sb, "+ %-30s (added)\n", e.Key)
		case Removed:
			fmt.Fprintf(&sb, "- %-30s (removed)\n", e.Key)
		case Modified:
			fmt.Fprintf(&sb, "~ %-30s (modified)\n", e.Key)
		}
	}
	return sb.String()
}

// Summary returns a one-line summary of the audit report.
func Summary(r Report) string {
	if r.IsClean() {
		return "no changes"
	}
	var added, removed, modified int
	for _, e := range r.Entries {
		switch e.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("%d added, %d removed, %d modified", added, removed, modified)
}
