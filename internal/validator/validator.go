// Package validator checks .env values against simple rules such as
// required keys, non-empty values, and basic format patterns.
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule applied to a key/value pair.
type Rule struct {
	Key     string
	Pattern string // optional regex pattern the value must match
	Required bool
}

// Violation describes a single validation failure.
type Violation struct {
	Key     string
	Value   string
	Message string
}

// String returns a human-readable description of the violation.
func (v Violation) String() string {
	return fmt.Sprintf("key %q: %s", v.Key, v.Message)
}

// Validate applies the given rules to the provided env map and returns
// any violations found. An empty slice means the env is valid.
func Validate(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.Required && strings.TrimSpace(val) == "" {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Value:   val,
				Message: "required key has empty value",
			})
			continue
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Value:   val,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Value:   val,
					Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern),
				})
			}
		}
	}

	return violations
}
