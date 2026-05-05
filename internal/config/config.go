// Package config defines the CLI configuration and validation logic for envcmp.
package config

import "errors"

// Config holds all parsed CLI options for a comparison run.
type Config struct {
	LeftFile   string
	RightFile  string
	Format     string
	Strict     bool
	IgnoreFile string

	// Filter options
	OnlyMissing    bool
	OnlyMismatched bool
	Keys           []string

	// Sort output
	Sort bool
}

// DefaultFormat is used when no format flag is provided.
const DefaultFormat = "text"

// Validate checks that the config is complete and applies defaults.
func Validate(c *Config) error {
	if c.LeftFile == "" {
		return errors.New("left file is required")
	}
	if c.RightFile == "" {
		return errors.New("right file is required")
	}
	if c.Format == "" {
		c.Format = DefaultFormat
	}
	switch c.Format {
	case "text", "json", "markdown":
		// valid
	default:
		return errors.New("format must be one of: text, json, markdown")
	}
	if c.OnlyMissing && c.OnlyMismatched {
		return errors.New("--only-missing and --only-mismatched are mutually exclusive")
	}
	return nil
}
