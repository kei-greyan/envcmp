// Package envscope extracts a scoped subset of an environment map by
// matching keys against a prefix, suffix, or both.
//
// It is useful for isolating service-specific variables (e.g. all keys
// beginning with "APP_" or "DB_") from a larger shared env file, and
// optionally stripping the matched prefix or suffix from the resulting keys
// to produce a clean, namespaced map.
//
// Example:
//
//	r, err := envscope.Extract(env, envscope.Options{
//		Prefix:      "APP_",
//		StripPrefix: true,
//	})
package envscope
