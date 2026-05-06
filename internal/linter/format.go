package linter

import (
	"fmt"
	"strings"
)

// Format returns a human-readable string representation of a lint Result.
// Each issue is printed on its own line with the file name, line number, key,
// and message. If the result is clean an empty string is returned.
func Format(r Result) string {
	if r.IsClean() {
		return ""
	}

	var sb strings.Builder
	for _, iss := range r.Issues {
		sb.WriteString(fmt.Sprintf("%s:%d [%s] %s\n", r.File, iss.Line, iss.Key, iss.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Summary returns a one-line summary such as "3 issue(s) found in prod.env".
func Summary(r Result) string {
	if r.IsClean() {
		return fmt.Sprintf("no issues found in %s", r.File)
	}
	return fmt.Sprintf("%d issue(s) found in %s", len(r.Issues), r.File)
}
