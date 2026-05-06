// Package linter provides static analysis for individual .env files.
//
// It detects common authoring mistakes before a diff is performed:
//
//   - Duplicate keys within the same file
//   - Key names that do not follow the POSIX convention (A-Z, digits, underscore,
//     must not start with a digit)
//   - Empty values that may indicate an unset placeholder
//   - Lines that are neither comments, blank lines, nor valid KEY=VALUE pairs
//
// Usage:
//
//	result := linter.Lint(filename, rawLines, parsedEnvMap)
//	if !result.IsClean() {
//		for _, issue := range result.Issues {
//			fmt.Printf("line %d [%s]: %s\n", issue.Line, issue.Key, issue.Message)
//		}
//	}
package linter
