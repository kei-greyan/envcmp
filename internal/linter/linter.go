// Package linter checks individual .env files for common issues such as
// duplicate keys, invalid key names, and empty values.
package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Issue represents a single linting problem found in an env file.
type Issue struct {
	Line    int
	Key     string
	Message string
}

// Result holds all issues found during a lint pass.
type Result struct {
	File   string
	Issues []Issue
}

// IsClean returns true when no issues were found.
func (r Result) IsClean() bool {
	return len(r.Issues) == 0
}

var validKeyRe = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// Lint analyses the parsed key-value map produced by parser.ParseFile together
// with the raw lines of the file so that line numbers can be reported.
func Lint(file string, lines []string, env map[string]string) Result {
	result := Result{File: file}
	seen := make(map[string]int)

	for i, raw := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) < 2 {
			result.Issues = append(result.Issues, Issue{
				Line:    lineNum,
				Key:     trimmed,
				Message: "line is not a valid KEY=VALUE assignment",
			})
			continue
		}

		key := strings.TrimSpace(parts[0])

		if !validKeyRe.MatchString(key) {
			result.Issues = append(result.Issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: fmt.Sprintf("key %q does not match [A-Z_][A-Z0-9_]*", key),
			})
		}

		if prev, dup := seen[key]; dup {
			result.Issues = append(result.Issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: fmt.Sprintf("duplicate key (first seen on line %d)", prev),
			})
		} else {
			seen[key] = lineNum
		}

		val := strings.TrimSpace(parts[1])
		if val == "" {
			result.Issues = append(result.Issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: "value is empty",
			})
		}
	}

	return result
}
