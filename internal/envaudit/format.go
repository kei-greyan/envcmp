package envaudit

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenderJSON returns the audit report serialised as indented JSON.
func RenderJSON(r Report) (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("envaudit: json render: %w", err)
	}
	return string(b), nil
}

// RenderMarkdown returns the audit report as a Markdown table.
func RenderMarkdown(r Report) string {
	if r.IsClean() {
		return "_No changes detected._\n"
	}
	var sb strings.Builder
	sb.WriteString("| Key | Change | Old Value | New Value |\n")
	sb.WriteString("|-----|--------|-----------|-----------|\n")
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "| %s | %s | %s | %s |\n",
			e.Key,
			string(e.Kind),
			maskIfEmpty(e.OldValue),
			maskIfEmpty(e.NewValue),
		)
	}
	return sb.String()
}

func maskIfEmpty(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
