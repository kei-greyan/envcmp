// Package envprune removes keys from an env map based on configurable criteria.
package envprune

import "strings"

// Options controls which keys are pruned.
type Options struct {
	// RemoveEmpty removes keys whose value is an empty string.
	RemoveEmpty bool
	// RemovePrefixes removes keys that start with any of the given prefixes.
	RemovePrefixes []string
	// RemoveSuffixes removes keys that end with any of the given suffixes.
	RemoveSuffixes []string
	// RemoveKeys removes exactly the listed key names.
	RemoveKeys []string
}

// Result holds the output of a prune operation.
type Result struct {
	Kept    map[string]string
	Pruned  map[string]string
}

// Prune applies the given options to env and returns a Result.
func Prune(env map[string]string, opts Options) Result {
	exact := make(map[string]struct{}, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		exact[k] = struct{}{}
	}

	kept := make(map[string]string)
	pruned := make(map[string]string)

	for k, v := range env {
		if shouldPrune(k, v, opts, exact) {
			pruned[k] = v
		} else {
			kept[k] = v
		}
	}

	return Result{Kept: kept, Pruned: pruned}
}

func shouldPrune(key, value string, opts Options, exact map[string]struct{}) bool {
	if _, ok := exact[key]; ok {
		return true
	}
	if opts.RemoveEmpty && value == "" {
		return true
	}
	for _, p := range opts.RemovePrefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	for _, s := range opts.RemoveSuffixes {
		if strings.HasSuffix(key, s) {
			return true
		}
	}
	return false
}
