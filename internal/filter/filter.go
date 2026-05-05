package filter

import "github.com/user/envcmp/internal/comparator"

// Options holds filtering criteria for diff results.
type Options struct {
	// OnlyMissing limits results to keys missing in either side.
	OnlyMissing bool
	// OnlyMismatched limits results to keys present in both sides but with different values.
	OnlyMismatched bool
	// Keys restricts results to a specific set of key names (empty means all).
	Keys []string
}

// Apply filters a comparator.Result according to the provided Options and
// returns a new Result containing only the matching entries.
func Apply(result comparator.Result, opts Options) comparator.Result {
	keySet := buildKeySet(opts.Keys)

	out := comparator.Result{
		MissingInRight: []string{},
		MissingInLeft:  []string{},
		Mismatched:     []comparator.Diff{},
	}

	if !opts.OnlyMismatched {
		for _, k := range result.MissingInRight {
			if matchesKeySet(k, keySet) {
				out.MissingInRight = append(out.MissingInRight, k)
			}
		}
		for _, k := range result.MissingInLeft {
			if matchesKeySet(k, keySet) {
				out.MissingInLeft = append(out.MissingInLeft, k)
			}
		}
	}

	if !opts.OnlyMissing {
		for _, d := range result.Mismatched {
			if matchesKeySet(d.Key, keySet) {
				out.Mismatched = append(out.Mismatched, d)
			}
		}
	}

	return out
}

// IsEmpty reports whether a Result contains no differences after filtering.
// This is useful for callers that want to check if two environments are
// equivalent under the given filter options without inspecting each field.
func IsEmpty(result comparator.Result) bool {
	return len(result.MissingInRight) == 0 &&
		len(result.MissingInLeft) == 0 &&
		len(result.Mismatched) == 0
}

// buildKeySet converts a slice of key names into a lookup map.
// An empty slice signals "no filter" (all keys match).
func buildKeySet(keys []string) map[string]struct{} {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

// matchesKeySet returns true when the key should be included in the output.
func matchesKeySet(key string, keySet map[string]struct{}) bool {
	if keySet == nil {
		return true
	}
	_, ok := keySet[key]
	return ok
}
