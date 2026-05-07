// Package envcheck verifies that all keys declared in a reference
// environment file (e.g. .env.example) are present — and optionally
// non-empty — in a target environment map.
//
// Typical usage:
//
//	reference, _ := parser.ParseFile(".env.example")
//	live, _      := parser.ParseFile(".env")
//	result := envcheck.Check(reference, live, envcheck.Options{FailOnEmpty: true})
//	if !result.IsClean() {
//		fmt.Println(envcheck.Format(result))
//	}
package envcheck
