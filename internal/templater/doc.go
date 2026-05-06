// Package templater provides utilities for generating .env template files
// from a parsed environment variable map.
//
// A template file contains all keys from the source env file with their
// values stripped (replaced by an empty string or a custom placeholder).
// This is useful for sharing env structure without exposing sensitive values.
//
// Example usage:
//
//	env, _ := parser.ParseFile(".env.production")
//	templater.WriteFile(".env.template", env, templater.Options{
//		Placeholder: "CHANGEME",
//		AddComments: true,
//	})
package templater
