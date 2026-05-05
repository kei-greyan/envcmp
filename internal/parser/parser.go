package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env represents a parsed .env file as a map of key-value pairs.
type Env map[string]string

// ParseFile reads a .env file from the given path and returns an Env map.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are skipped. Keys without values are stored with an empty string.
func ParseFile(path string) (Env, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: could not open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(Env)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("parser: empty key on line %d in %q", lineNum, path)
		}

		value := ""
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
			value = stripQuotes(value)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: error reading %q: %w", path, err)
	}

	return env, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
