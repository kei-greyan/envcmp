// Package envresolve expands variable references within an env map.
//
// Values may reference other keys using ${VAR} or $VAR syntax. Resolution
// is performed iteratively up to a configurable depth to handle chained
// references while guarding against infinite cycles.
//
// Example:
//
//	env := map[string]string{
//		"HOST":   "localhost",
//		"PORT":   "5432",
//		"DB_URL": "postgres://${HOST}:${PORT}/mydb",
//	}
//	result, err := envresolve.Resolve(env, envresolve.DefaultOptions())
//	// result.Resolved["DB_URL"] == "postgres://localhost:5432/mydb"
package envresolve
