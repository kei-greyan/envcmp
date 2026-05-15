// Package envgroup partitions a flat map of environment variables into named
// groups based on a configurable key delimiter and depth.
//
// Example:
//
//	env := map[string]string{
//		"DB_HOST": "localhost",
//		"DB_PORT": "5432",
//		"APP_ENV":  "production",
//	}
//	 result, _ := envgroup.Group(env, envgroup.DefaultOptions())
//	// Groups: {"APP": {...}, "DB": {...}}
package envgroup
