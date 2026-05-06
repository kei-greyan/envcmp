// Package redactor masks sensitive values in comparison results
// before they are passed to formatters or reporters.
package redactor

import (
	"strings"

	"github.com/user/envcmp/internal/comparator"
)

const masked = "***"

// sensitivePatterns holds substrings that indicate a key is sensitive.
var sensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
}

// IsSensitive reports whether the given key name looks sensitive.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// Apply returns a copy of result with sensitive values replaced by "***".
// Keys are never masked — only their values.
func Apply(result comparator.Result) comparator.Result {
	out := comparator.Result{
		MissingInRight: make([]string, len(result.MissingInRight)),
		MissingInLeft:  make([]string, len(result.MissingInLeft)),
		Mismatched:     make([]comparator.Diff, len(result.Mismatched)),
	}

	copy(out.MissingInRight, result.MissingInRight)
	copy(out.MissingInLeft, result.MissingInLeft)

	for i, d := range result.Mismatched {
		entry := comparator.Diff{Key: d.Key}
		if IsSensitive(d.Key) {
			entry.LeftVal = masked
			entry.RightVal = masked
		} else {
			entry.LeftVal = d.LeftVal
			entry.RightVal = d.RightVal
		}
		out.Mismatched[i] = entry
	}

	return out
}
