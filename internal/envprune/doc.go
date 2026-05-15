// Package envprune provides utilities for pruning keys from an environment map.
//
// Keys can be removed by:
//   - exact name match
//   - empty value
//   - key prefix
//   - key suffix
//
// The Prune function returns a Result containing the kept and pruned entries,
// allowing callers to inspect what was removed without mutating the original map.
package envprune
