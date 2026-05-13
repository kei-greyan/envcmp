// Package envnorm normalizes environment variable maps before comparison or
// export. Normalization includes trimming whitespace from values and
// optionally converting key casing to a canonical form.
//
// Usage:
//
//	opts := envnorm.DefaultOptions()
//	normalized := envnorm.Normalize(env, opts)
//
// Use Diff to preview which keys or values would be affected without
// committing to the full normalization:
//
//	changed := envnorm.Diff(env, opts)
package envnorm
