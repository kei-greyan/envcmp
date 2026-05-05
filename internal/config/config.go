// Package config handles loading and validation of envcmp runtime configuration
// from flags, environment variables, or a config file.
package config

import (
	"errors"
	"strings"
)

// Format represents the supported output formats.
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
)

// Config holds the resolved runtime configuration for a single envcmp run.
type Config struct {
	// LeftFile is the path to the base .env file.
	LeftFile string
	// RightFile is the path to the comparison .env file.
	RightFile string
	// Format controls the output rendering format.
	Format Format
	// Strict causes the process to exit non-zero when any diff is found.
	Strict bool
	// OnlyMissing limits output to keys missing in either file.
	OnlyMissing bool
	// OnlyMismatched limits output to keys present in both files but with different values.
	OnlyMismatched bool
	// Keys is an optional allowlist of specific keys to compare.
	Keys []string
}

// Validate returns an error if the configuration is not usable.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.LeftFile) == "" {
		return errors.New("left file path must not be empty")
	}
	if strings.TrimSpace(c.RightFile) == "" {
		return errors.New("right file path must not be empty")
	}
	switch c.Format {
	case FormatText, FormatJSON, FormatMarkdown:
		// valid
	case "":
		c.Format = FormatText
	default:
		return errors.New("unsupported format: " + string(c.Format))
	}
	if c.OnlyMissing && c.OnlyMismatched {
		return errors.New("--only-missing and --only-mismatched are mutually exclusive")
	}
	return nil
}
