// Package envaudit compares two snapshots of an environment variable map
// and produces a structured audit Report describing every addition, removal,
// and modification.
//
// Basic usage:
//
//	before := map[string]string{"APP_ENV": "production", "DB_HOST": "db1"}
//	after  := map[string]string{"APP_ENV": "staging",    "DB_HOST": "db1", "NEW_KEY": "x"}
//
//	report := envaudit.Audit(before, after)
//	fmt.Print(envaudit.Format(report))
//	fmt.Println(envaudit.Summary(report))
//
// For machine-readable output use RenderJSON or RenderMarkdown.
package envaudit
