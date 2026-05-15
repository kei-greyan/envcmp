// Package envtag provides tagging and grouping of env keys by metadata labels.
package envtag

import (
	"sort"
	"strings"
)

// Tag represents a named label applied to a set of env keys.
type Tag struct {
	Name string
	Keys []string
}

// Result holds the outcome of a tagging operation.
type Result struct {
	Tagged   map[string][]string // tag name -> list of keys
	Untagged []string            // keys with no matching tag
}

// Options controls tagging behaviour.
type Options struct {
	// Tags maps tag names to key prefixes or exact keys.
	// A key matches a tag if it starts with one of the tag's prefixes.
	Tags map[string][]string
	ExactMatch bool
}

// Apply tags every key in env according to opts and returns a Result.
func Apply(env map[string]string, opts Options) Result {
	result := Result{
		Tagged: make(map[string][]string),
	}

	for key := range env {
		matched := false
		for tagName, patterns := range opts.Tags {
			if matchesTag(key, patterns, opts.ExactMatch) {
				result.Tagged[tagName] = append(result.Tagged[tagName], key)
				matched = true
			}
		}
		if !matched {
			result.Untagged = append(result.Untagged, key)
		}
	}

	for tagName := range result.Tagged {
		sort.Strings(result.Tagged[tagName])
	}
	sort.Strings(result.Untagged)

	return result
}

// Keys returns all keys associated with a given tag name, or nil if not found.
func Keys(r Result, tag string) []string {
	return r.Tagged[tag]
}

// TagNames returns a sorted list of all tag names present in the result.
func TagNames(r Result) []string {
	names := make([]string, 0, len(r.Tagged))
	for name := range r.Tagged {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func matchesTag(key string, patterns []string, exact bool) bool {
	for _, p := range patterns {
		if exact {
			if key == p {
				return true
			}
		} else {
			if strings.HasPrefix(key, p) {
				return true
			}
		}
	}
	return false
}
