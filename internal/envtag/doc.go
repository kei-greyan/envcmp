// Package envtag provides utilities for tagging and grouping environment
// variable keys by user-defined labels.
//
// Tags are defined as mappings from a label name to a list of key prefixes
// (or exact key names when ExactMatch is enabled). Each key in the env map
// is evaluated against every tag and assigned accordingly. Keys that do not
// match any tag are placed in the Untagged list.
//
// Example usage:
//
//	opts := envtag.Options{
//		Tags: map[string][]string{
//			"database": {"DB_", "POSTGRES_"},
//			"auth":     {"AUTH_", "JWT_"},
//		},
//	}
//	result := envtag.Apply(env, opts)
//	fmt.Println(envtag.Format(result))
package envtag
