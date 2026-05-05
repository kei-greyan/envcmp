// Package ignore provides functionality to load and apply ignore rules
// for env keys that should be excluded from comparison results.
package ignore

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envcmp/internal/comparator"
)

// LoadFile reads an ignore file and returns a set of key patterns to ignore.
// Each non-empty, non-comment line is treated as a key name to ignore.
func LoadFile(path string) (map[string]struct{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	keys := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		keys[line] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

// Apply removes any entries from result whose key appears in the ignore set.
func Apply(result comparator.Result, ignoreKeys map[string]struct{}) comparator.Result {
	if len(ignoreKeys) == 0 {
		return result
	}

	filtered := comparator.Result{}

	for _, k := range result.MissingInRight {
		if _, skip := ignoreKeys[k]; !skip {
			filtered.MissingInRight = append(filtered.MissingInRight, k)
		}
	}
	for _, k := range result.MissingInLeft {
		if _, skip := ignoreKeys[k]; !skip {
			filtered.MissingInLeft = append(filtered.MissingInLeft, k)
		}
	}
	for _, mm := range result.Mismatched {
		if _, skip := ignoreKeys[mm.Key]; !skip {
			filtered.Mismatched = append(filtered.Mismatched, mm)
		}
	}

	return filtered
}
