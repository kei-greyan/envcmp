// Package schema provides functionality for validating .env files
// against a declared schema of expected keys and their constraints.
package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// FieldType represents the expected type of a value.
type FieldType string

const (
	TypeString FieldType = "string"
	TypeInt    FieldType = "int"
	TypeBool   FieldType = "bool"
	TypeURL    FieldType = "url"
)

// Field describes constraints for a single key.
type Field struct {
	Required bool      `json:"required"`
	Type     FieldType `json:"type"`
	Pattern  string    `json:"pattern"`
}

// Schema maps key names to their field definitions.
type Schema map[string]Field

// Violation describes a single schema violation.
type Violation struct {
	Key     string
	Message string
}

var (
	reInt  = regexp.MustCompile(`^-?\d+$`)
	reBool = regexp.MustCompile(`^(true|false|1|0|yes|no)$`)
	reURL  = regexp.MustCompile(`^https?://`)
)

// LoadFile reads and parses a JSON schema file.
func LoadFile(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read %q: %w", path, err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("schema: parse %q: %w", path, err)
	}
	return s, nil
}

// Validate checks env against the schema and returns any violations.
func Validate(env map[string]string, s Schema) []Violation {
	var violations []Violation
	for key, field := range s {
		val, exists := env[key]
		if !exists || val == "" {
			if field.Required {
				violations = append(violations, Violation{Key: key, Message: "required key is missing or empty"})
			}
			continue
		}
		if v := validateType(key, val, field.Type); v != nil {
			violations = append(violations, *v)
		}
		if field.Pattern != "" {
			if v := validatePattern(key, val, field.Pattern); v != nil {
				violations = append(violations, *v)
			}
		}
	}
	return violations
}

func validateType(key, val string, t FieldType) *Violation {
	switch t {
	case TypeInt:
		if !reInt.MatchString(val) {
			return &Violation{Key: key, Message: fmt.Sprintf("expected int, got %q", val)}
		}
	case TypeBool:
		if !reBool.MatchString(val) {
			return &Violation{Key: key, Message: fmt.Sprintf("expected bool, got %q", val)}
		}
	case TypeURL:
		if !reURL.MatchString(val) {
			return &Violation{Key: key, Message: fmt.Sprintf("expected URL, got %q", val)}
		}
	}
	return nil
}

func validatePattern(key, val, pattern string) *Violation {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return &Violation{Key: key, Message: fmt.Sprintf("invalid pattern %q: %v", pattern, err)}
	}
	if !re.MatchString(val) {
		return &Violation{Key: key, Message: fmt.Sprintf("value %q does not match pattern %q", val, pattern)}
	}
	return nil
}
