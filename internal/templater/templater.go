// Package templater generates a .env.template file from a parsed env map,
// replacing all values with empty strings or placeholder comments.
package templater

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Options controls how the template is generated.
type Options struct {
	// Placeholder is written as the value for each key. Defaults to "".
	Placeholder string
	// AddComments prepends a comment above each key when true.
	AddComments bool
}

// Generate builds a template string from the provided env map.
// Keys are sorted alphabetically for deterministic output.
func Generate(env map[string]string, opts Options) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		if opts.AddComments {
			fmt.Fprintf(&sb, "# %s\n", k)
		}
		fmt.Fprintf(&sb, "%s=%s\n", k, opts.Placeholder)
	}
	return sb.String()
}

// WriteFile writes the generated template to the given file path.
// Returns an error if the file cannot be created or written.
func WriteFile(path string, env map[string]string, opts Options) error {
	content := Generate(env, opts)
	return os.WriteFile(path, []byte(content), 0o644)
}
