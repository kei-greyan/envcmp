// Package envcast provides utilities for casting environment variable
// string values to typed Go values with validation support.
package envcast

import (
	"fmt"
	"strconv"
	"strings"
)

// Result holds the outcome of a cast operation for a single key.
type Result struct {
	Key   string
	Raw   string
	Kind  string
	Value any
	Err   error
}

// CastOptions controls how values are cast.
type CastOptions struct {
	// Strict causes Cast to return an error on the first failure.
	Strict bool
}

// Cast attempts to convert each key in env to the type specified in types.
// types maps a key name to one of: "string", "int", "float", "bool".
// Keys not present in types are returned as-is with kind "string".
func Cast(env map[string]string, types map[string]string, opts CastOptions) ([]Result, error) {
	results := make([]Result, 0, len(env))

	for key, raw := range env {
		kind, ok := types[key]
		if !ok {
			kind = "string"
		}

		r := Result{Key: key, Raw: raw, Kind: kind}
		var err error

		switch strings.ToLower(kind) {
		case "string":
			r.Value = raw
		case "int":
			r.Value, err = strconv.Atoi(raw)
		case "float":
			r.Value, err = strconv.ParseFloat(raw, 64)
		case "bool":
			r.Value, err = strconv.ParseBool(raw)
		default:
			err = fmt.Errorf("unknown type %q", kind)
		}

		if err != nil {
			r.Err = fmt.Errorf("key %q: cannot cast %q to %s: %w", key, raw, kind, err)
			if opts.Strict {
				return nil, r.Err
			}
		}

		results = append(results, r)
	}

	return results, nil
}

// Failures returns only the results that contain a cast error.
func Failures(results []Result) []Result {
	out := make([]Result, 0)
	for _, r := range results {
		if r.Err != nil {
			out = append(out, r)
		}
	}
	return out
}
