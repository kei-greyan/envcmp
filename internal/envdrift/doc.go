// Package envdrift detects configuration drift between a persisted baseline
// snapshot and a live environment map.
//
// Use Detect to compare two maps and obtain a Result containing all
// added, removed, and changed keys. Use Format to render a human-readable
// diff and Summary for a concise one-line status message.
//
// Typical usage:
//
//	baseline, _ := envpin.LoadFile("baseline.json")
//	current, _ := parser.ParseFile(".env")
//	result := envdrift.Detect(baseline, current)
//	fmt.Print(envdrift.Format(result))
package envdrift
