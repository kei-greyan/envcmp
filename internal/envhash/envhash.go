// Package envhash provides utilities for computing and comparing
// stable content hashes of environment variable maps.
package envhash

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Result holds the hash output for an environment map.
type Result struct {
	Hash    string            // hex-encoded SHA-256 of the canonical form
	Entries map[string]string // per-key hashes
}

// Compute returns a Result for the given env map.
// The top-level hash is derived from a stable, sorted canonical representation
// so that key order does not affect the output.
func Compute(env map[string]string) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make(map[string]string, len(env))
	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		entry := fmt.Sprintf("%s=%s", k, v)
		h := sha256.Sum256([]byte(entry))
		entries[k] = fmt.Sprintf("%x", h)
		fmt.Fprintf(&sb, "%s\n", entry)
	}

	total := sha256.Sum256([]byte(sb.String()))
	return Result{
		Hash:    fmt.Sprintf("%x", total),
		Entries: entries,
	}
}

// Equal reports whether two Results share the same top-level hash.
func Equal(a, b Result) bool {
	return a.Hash == b.Hash
}

// Diff returns the keys whose per-entry hashes differ between a and b.
// Keys present in one result but not the other are also included.
func Diff(a, b Result) []string {
	seen := make(map[string]struct{})
	for k := range a.Entries {
		seen[k] = struct{}{}
	}
	for k := range b.Entries {
		seen[k] = struct{}{}
	}

	var changed []string
	for k := range seen {
		if a.Entries[k] != b.Entries[k] {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	return changed
}
