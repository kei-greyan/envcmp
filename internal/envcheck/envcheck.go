// Package envcheck provides functionality to verify that all keys
// defined in a reference .env file are present in the current environment.
package envcheck

import "fmt"

// Result holds the outcome of an environment check.
type Result struct {
	Present []string
	Missing []string
	Empty   []string
}

// IsClean returns true if there are no missing or empty keys.
func (r Result) IsClean() bool {
	return len(r.Missing) == 0 && len(r.Empty) == 0
}

// Options controls the behaviour of Check.
type Options struct {
	// FailOnEmpty treats keys that exist but have empty values as failures.
	FailOnEmpty bool
}

// Check verifies that every key in reference exists in env.
// reference is typically parsed from a .env.example file.
// env is the live environment map (e.g. from os.Environ or a parsed .env).
func Check(reference map[string]string, env map[string]string, opts Options) Result {
	result := Result{}

	for key := range reference {
		val, exists := env[key]
		if !exists {
			result.Missing = append(result.Missing, key)
			continue
		}
		if opts.FailOnEmpty && val == "" {
			result.Empty = append(result.Empty, key)
			continue
		}
		result.Present = append(result.Present, key)
	}

	return result
}

// Format returns a human-readable report of the check result.
func Format(r Result) string {
	if r.IsClean() {
		return fmt.Sprintf("OK: all %d keys present", len(r.Present))
	}
	out := ""
	for _, k := range r.Missing {
		out += fmt.Sprintf("MISSING   %s\n", k)
	}
	for _, k := range r.Empty {
		out += fmt.Sprintf("EMPTY     %s\n", k)
	}
	return out
}

// Summary returns a one-line summary suitable for CLI output.
func Summary(r Result) string {
	return fmt.Sprintf("%d present, %d missing, %d empty",
		len(r.Present), len(r.Missing), len(r.Empty))
}
