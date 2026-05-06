// Package schema provides JSON-based schema validation for .env files.
//
// A schema is a JSON file mapping environment variable names to field
// definitions that declare whether a key is required, its expected type
// (string, int, bool, url), and an optional regex pattern the value must
// satisfy.
//
// Example schema.json:
//
//	{
//	  "PORT":    { "required": true, "type": "int" },
//	  "API_URL": { "required": true, "type": "url" },
//	  "APP_ENV": { "required": true, "pattern": "^(dev|staging|prod)$" }
//	}
//
// Use LoadFile to parse a schema file, then Validate to check an env map.
package schema
