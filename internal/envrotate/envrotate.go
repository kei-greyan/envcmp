// Package envrotate provides utilities for rotating values in .env files,
// generating new secrets or placeholders for specified keys while preserving
// the rest of the environment map.
package envrotate

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
)

// Strategy controls how new values are generated during rotation.
type Strategy string

const (
	// StrategyRandom replaces the value with a cryptographically random hex string.
	StrategyRandom Strategy = "random"
	// StrategyBlank replaces the value with an empty string.
	StrategyBlank Strategy = "blank"
	// StrategyPlaceholder replaces the value with a descriptive placeholder.
	StrategyPlaceholder Strategy = "placeholder"
)

// Entry records the old and new value for a rotated key.
type Entry struct {
	Key      string
	OldValue string
	NewValue string
}

// Result holds the updated environment map and the rotation log.
type Result struct {
	Env     map[string]string
	Rotated []Entry
}

// Options configures the rotation behaviour.
type Options struct {
	Keys     []string
	Strategy Strategy
	// ByteLen is the number of random bytes used when Strategy is StrategyRandom.
	ByteLen int
}

// DefaultOptions returns sensible defaults for rotation.
func DefaultOptions() Options {
	return Options{
		Strategy: StrategyRandom,
		ByteLen:  16,
	}
}

// Rotate applies the rotation to the provided environment map and returns a
// Result containing the updated map and a log of every change made.
// The original map is never mutated.
func Rotate(env map[string]string, opts Options) (Result, error) {
	if len(opts.Keys) == 0 {
		return Result{}, fmt.Errorf("envrotate: no keys specified for rotation")
	}
	if opts.ByteLen <= 0 {
		opts.ByteLen = 16
	}

	updated := make(map[string]string, len(env))
	for k, v := range env {
		updated[k] = v
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	rotated := make([]Entry, 0, len(opts.Keys))
	sorted := make([]string, 0, len(opts.Keys))
	for _, k := range opts.Keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, key := range sorted {
		oldVal := updated[key]
		newVal, err := generateValue(key, opts)
		if err != nil {
			return Result{}, fmt.Errorf("envrotate: failed to generate value for %q: %w", key, err)
		}
		updated[key] = newVal
		rotated = append(rotated, Entry{Key: key, OldValue: oldVal, NewValue: newVal})
	}

	return Result{Env: updated, Rotated: rotated}, nil
}

func generateValue(key string, opts Options) (string, error) {
	switch opts.Strategy {
	case StrategyRandom:
		b := make([]byte, opts.ByteLen)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return hex.EncodeToString(b), nil
	case StrategyBlank:
		return "", nil
	case StrategyPlaceholder:
		return fmt.Sprintf("ROTATED_%s", key), nil
	default:
		return "", fmt.Errorf("unknown strategy %q", opts.Strategy)
	}
}
