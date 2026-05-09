// Package envcast provides type-casting utilities for environment variable values.
//
// It converts raw string values from a parsed .env map into typed Go values
// such as int, float64, bool, or string. Casting is driven by a type-map
// supplied by the caller, making it easy to integrate with schema or validator
// workflows.
//
// Usage:
//
//	types := map[string]string{"PORT": "int", "DEBUG": "bool"}
//	results, err := envcast.Cast(env, types, envcast.CastOptions{Strict: true})
//	fails := envcast.Failures(results)
package envcast
