// Package envpin records a snapshot of expected env keys and their types,
// then checks a live env map against that pin to detect drift.
package envpin

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// Pin is a recorded snapshot of env key metadata.
type Pin struct {
	Keys map[string]PinEntry `json:"keys"`
}

// PinEntry holds metadata about a single key at pin time.
type PinEntry struct {
	Present  bool   `json:"present"`
	NonEmpty bool   `json:"non_empty"`
	Kind     string `json:"kind"` // "string", "int", "bool", "float"
}

// Violation describes a drift detected against the pin.
type Violation struct {
	Key     string
	Message string
}

// Create builds a Pin from the provided env map.
func Create(env map[string]string) Pin {
	p := Pin{Keys: make(map[string]PinEntry, len(env))}
	for k, v := range env {
		p.Keys[k] = PinEntry{
			Present:  true,
			NonEmpty: v != "",
			Kind:     inferKind(v),
		}
	}
	return p
}

// SaveFile writes a Pin to a JSON file at path.
func SaveFile(path string, p Pin) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("envpin: create %q: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

// LoadFile reads a Pin from a JSON file at path.
func LoadFile(path string) (Pin, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Pin{}, fmt.Errorf("envpin: read %q: %w", path, err)
	}
	var p Pin
	if err := json.Unmarshal(data, &p); err != nil {
		return Pin{}, fmt.Errorf("envpin: parse %q: %w", path, err)
	}
	return p, nil
}

// Check compares env against pin and returns any violations.
func Check(p Pin, env map[string]string) []Violation {
	var violations []Violation
	keys := make([]string, 0, len(p.Keys))
	for k := range p.Keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		entry := p.Keys[k]
		v, exists := env[k]
		if entry.Present && !exists {
			violations = append(violations, Violation{Key: k, Message: "key missing from env"})
			continue
		}
		if entry.NonEmpty && v == "" {
			violations = append(violations, Violation{Key: k, Message: "expected non-empty value"})
			continue
		}
		if exists && v != "" && entry.Kind != "string" {
			if got := inferKind(v); got != entry.Kind {
				violations = append(violations, Violation{
					Key:     k,
					Message: fmt.Sprintf("type drift: pinned %s, got %s", entry.Kind, got),
				})
			}
		}
	}
	return violations
}

func inferKind(v string) string {
	if _, err := strconv.ParseBool(v); err == nil {
		return "bool"
	}
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return "int"
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return "float"
	}
	return "string"
}
