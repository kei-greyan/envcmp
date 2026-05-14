// Package envmerge provides utilities for merging multiple .env maps
// into a single resolved map, with configurable conflict strategies
// and a detailed change log of what was merged or overridden.
package envmerge

import (
	"fmt"
	"sort"
)

// Strategy controls how key conflicts are resolved across sources.
type Strategy string

const (
	StrategyFirst Strategy = "first" // keep the first value seen
	StrategyLast  Strategy = "last"  // keep the last value seen
	StrategyError Strategy = "error" // return an error on conflict
)

// Entry records a single key's resolution during a merge.
type Entry struct {
	Key      string
	Value    string
	Source   int // index of the winning source
	Conflict bool
}

// Result holds the merged environment and the log of all entries.
type Result struct {
	Env     map[string]string
	Entries []Entry
}

// Merge combines multiple env maps using the given strategy.
// Sources are applied left-to-right; index 0 is the leftmost.
func Merge(strategy Strategy, sources ...map[string]string) (Result, error) {
	env := make(map[string]string)
	origin := make(map[string]int)
	conflicts := make(map[string]bool)

	for idx, src := range sources {
		for k, v := range src {
			if existing, found := env[k]; found {
				if existing == v {
					continue
				}
				switch strategy {
				case StrategyError:
					return Result{}, fmt.Errorf("conflict on key %q: %q vs %q", k, existing, v)
				case StrategyLast:
					env[k] = v
					origin[k] = idx
					conflicts[k] = true
				default: // StrategyFirst
					conflicts[k] = true
				}
			} else {
				env[k] = v
				origin[k] = idx
			}
		}
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{
			Key:      k,
			Value:    env[k],
			Source:   origin[k],
			Conflict: conflicts[k],
		})
	}

	return Result{Env: env, Entries: entries}, nil
}
