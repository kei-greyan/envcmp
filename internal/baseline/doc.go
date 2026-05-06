// Package baseline provides save and load functionality for persisting
// a comparator.Result as a named baseline file.
//
// A baseline captures the state of a .env comparison at a point in time,
// allowing subsequent runs to detect drift — new missing keys or value
// changes that were not present when the baseline was recorded.
//
// Usage:
//
//	if err := baseline.Save("baseline.json", left, right, result); err != nil { ... }
//
//	rec, err := baseline.Load("baseline.json")
package baseline
