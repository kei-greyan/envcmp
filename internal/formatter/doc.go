// Package formatter provides output rendering for envcmp comparison results.
//
// Supported formats:
//   - text     — human-readable plain text (default)
//   - json     — structured JSON suitable for machine consumption
//   - markdown — GitHub-flavoured markdown for use in PR comments or reports
//
// Usage:
//
//	out, err := formatter.Render(result, formatter.FormatJSON)
package formatter
