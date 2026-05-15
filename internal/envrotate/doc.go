// Package envrotate implements key rotation for .env environments.
//
// It accepts a map of key/value pairs and a list of keys to rotate, then
// generates new values according to a chosen Strategy:
//
//   - StrategyRandom      — cryptographically random hex string (default)
//   - StrategyBlank       — empty string (clears the value)
//   - StrategyPlaceholder — a descriptive placeholder derived from the key name
//
// The original environment map is never mutated; Rotate always returns a
// fresh copy together with a structured log of every change made.
package envrotate
