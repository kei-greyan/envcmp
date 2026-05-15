// Package envresolve resolves variable references within an env map,
// expanding values that reference other keys using ${VAR} or $VAR syntax.
package envresolve

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls resolution behaviour.
type Options struct {
	// MaxDepth limits recursive expansion to prevent cycles. Defaults to 10.
	MaxDepth int
	// FailOnMissing returns an error when a referenced key is absent.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{MaxDepth: 10}
}

// Result holds the resolved env map and any warnings produced during resolution.
type Result struct {
	Resolved map[string]string
	Warnings []string
}

// Resolve expands variable references in env values and returns a Result.
func Resolve(env map[string]string, opts Options) (Result, error) {
	if opts.MaxDepth <= 0 {
		opts.MaxDepth = 10
	}

	resolved := make(map[string]string, len(env))
	for k, v := range env {
		resolved[k] = v
	}

	var warnings []string

	for key := range resolved {
		val, warns, err := expand(key, resolved, opts, 0)
		if err != nil {
			return Result{}, err
		}
		resolved[key] = val
		warnings = append(warnings, warns...)
	}

	return Result{Resolved: resolved, Warnings: warnings}, nil
}

func expand(key string, env map[string]string, opts Options, depth int) (string, []string, error) {
	if depth > opts.MaxDepth {
		return env[key], nil, fmt.Errorf("envresolve: max expansion depth exceeded for key %q", key)
	}

	var warnings []string
	value := env[key]

	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		submatches := refPattern.FindStringSubmatch(match)
		refKey := submatches[1]
		if refKey == "" {
			refKey = submatches[2]
		}
		refVal, ok := env[refKey]
		if !ok {
			warnings = append(warnings, fmt.Sprintf("key %q references undefined variable %q", key, refKey))
			if opts.FailOnMissing {
				return match
			}
			return ""
		}
		if strings.Contains(refVal, "$") {
			expanded, _, _ := expand(refKey, env, opts, depth+1)
			return expanded
		}
		return refVal
	})

	if opts.FailOnMissing && len(warnings) > 0 {
		return "", warnings, fmt.Errorf("envresolve: unresolved reference in key %q", key)
	}

	return result, warnings, nil
}
