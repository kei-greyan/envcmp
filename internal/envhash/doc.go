// Package envhash computes stable SHA-256 content hashes for environment
// variable maps.
//
// The canonical form used for hashing sorts keys lexicographically and
// serialises each entry as "KEY=VALUE\n", ensuring that the resulting hash
// is deterministic regardless of map iteration order.
//
// Use Compute to hash a single env map, Equal to compare two hashes, and
// Diff to identify which specific keys changed between two snapshots.
package envhash
