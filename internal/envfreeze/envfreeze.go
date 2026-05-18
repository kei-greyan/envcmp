// Package envfreeze provides functionality to freeze an env map into a
// read-only snapshot and detect any mutations against that snapshot.
package envfreeze

import (
	"fmt"
	"sort"
)

// Snapshot is an immutable copy of an env map taken at freeze time.
type Snapshot struct {
	data map[string]string
}

// Freeze creates a Snapshot from the provided env map.
func Freeze(env map[string]string) Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Snapshot{data: copy}
}

// Get returns the value for a key in the snapshot and whether it existed.
func (s Snapshot) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

// Keys returns all keys in the snapshot in sorted order.
func (s Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Violation describes a single mutation detected against a frozen snapshot.
type Violation struct {
	Key      string
	Frozen   string
	Current  string
	Kind     string // "modified", "added", "removed"
}

// Check compares the current env map against the snapshot and returns any
// violations. An empty slice means the env is unchanged.
func Check(s Snapshot, current map[string]string) []Violation {
	var violations []Violation

	for k, frozen := range s.data {
		cur, ok := current[k]
		if !ok {
			violations = append(violations, Violation{Key: k, Frozen: frozen, Current: "", Kind: "removed"})
		} else if cur != frozen {
			violations = append(violations, Violation{Key: k, Frozen: frozen, Current: cur, Kind: "modified"})
		}
	}

	for k, cur := range current {
		if _, ok := s.data[k]; !ok {
			violations = append(violations, Violation{Key: k, Frozen: "", Current: cur, Kind: "added"})
		}
	}

	sort.Slice(violations, func(i, j int) bool {
		return violations[i].Key < violations[j].Key
	})
	return violations
}

// Summary returns a one-line human-readable summary of the check result.
func Summary(violations []Violation) string {
	if len(violations) == 0 {
		return "env is frozen: no mutations detected"
	}
	return fmt.Sprintf("env drift detected: %d mutation(s) against frozen snapshot", len(violations))
}
