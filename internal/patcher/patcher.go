// Package patcher applies a set of key-value overrides to an existing env map,
// producing a patched copy without mutating the original.
package patcher

import (
	"errors"
	"fmt"
	"strings"
)

// PatchMode controls how conflicts between the base env and patches are handled.
type PatchMode int

const (
	// Overwrite replaces existing keys with patch values.
	Overwrite PatchMode = iota
	// SkipExisting keeps existing values and ignores patch values for those keys.
	SkipExisting
	// ErrorOnConflict returns an error if a patch key already exists in the base env.
	ErrorOnConflict
)

// Result holds the patched environment map and metadata about what changed.
type Result struct {
	Env     map[string]string
	Added   []string
	Updated []string
	Skipped []string
}

// Apply merges patches into base according to mode and returns a Result.
// The original base map is never mutated.
func Apply(base map[string]string, patches map[string]string, mode PatchMode) (*Result, error) {
	if base == nil {
		return nil, errors.New("patcher: base env must not be nil")
	}
	if patches == nil {
		return nil, errors.New("patcher: patches map must not be nil")
	}

	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	result := &Result{Env: out}

	for k, v := range patches {
		if err := validateKey(k); err != nil {
			return nil, fmt.Errorf("patcher: invalid key %q: %w", k, err)
		}
		_, exists := out[k]
		switch {
		case !exists:
			out[k] = v
			result.Added = append(result.Added, k)
		case mode == Overwrite:
			out[k] = v
			result.Updated = append(result.Updated, k)
		case mode == SkipExisting:
			result.Skipped = append(result.Skipped, k)
		case mode == ErrorOnConflict:
			return nil, fmt.Errorf("patcher: conflict on key %q", k)
		}
	}

	return result, nil
}

// validateKey ensures the key is a non-empty, valid env var name.
func validateKey(k string) error {
	if strings.TrimSpace(k) == "" {
		return errors.New("key must not be empty or whitespace")
	}
	for i, ch := range k {
		if ch == '_' || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
			continue
		}
		if i > 0 && ch >= '0' && ch <= '9' {
			continue
		}
		return fmt.Errorf("invalid character %q at position %d", ch, i)
	}
	return nil
}
