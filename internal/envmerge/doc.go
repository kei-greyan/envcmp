// Package envmerge merges two or more parsed .env maps into a single
// resolved environment, applying a configurable conflict strategy.
//
// Supported strategies:
//
//   - first  – the first source that defines a key wins (default)
//   - last   – the last source that defines a key wins
//   - error  – any conflicting key returns an error immediately
//
// Identical values for the same key across sources are never treated
// as conflicts regardless of strategy.
//
// Usage:
//
//	r, err := envmerge.Merge(envmerge.StrategyFirst, mapA, mapB)
//	fmt.Println(envmerge.Summary(r))
package envmerge
