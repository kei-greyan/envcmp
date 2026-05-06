// Package differ provides utilities for computing and describing
// value-level differences between two sets of environment variables.
//
// Unlike the comparator package — which identifies missing keys —
// differ focuses on keys present in both environments whose values
// have changed. Results are returned as a sorted slice of ValueDiff
// structs that can be rendered or inspected programmatically.
//
// Typical usage:
//
//	diffs := differ.DiffValues(leftMap, rightMap)
//	fmt.Println(differ.Summary(diffs))
package differ
