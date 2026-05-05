package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envcmp/internal/comparator"
)

// Format defines the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatMarkdown Format = "markdown"
)

// JSONOutput is the structured representation for JSON format.
type JSONOutput struct {
	MissingInRight []string            `json:"missing_in_right,omitempty"`
	MissingInLeft  []string            `json:"missing_in_left,omitempty"`
	Mismatched     map[string][2]string `json:"mismatched,omitempty"`
}

// Render formats the comparison result according to the given format.
func Render(result comparator.Result, format Format) (string, error) {
	switch format {
	case FormatJSON:
		return renderJSON(result)
	case FormatMarkdown:
		return renderMarkdown(result), nil
	case FormatText:
		return renderText(result), nil
	default:
		return "", fmt.Errorf("unknown format: %q", format)
	}
}

func renderText(result comparator.Result) string {
	var sb strings.Builder
	for _, k := range result.MissingInRight {
		fmt.Fprintf(&sb, "[MISSING RIGHT] %s\n", k)
	}
	for _, k := range result.MissingInLeft {
		fmt.Fprintf(&sb, "[MISSING LEFT]  %s\n", k)
	}
	for k, v := range result.Mismatched {
		fmt.Fprintf(&sb, "[MISMATCH]      %s: %q != %q\n", k, v[0], v[1])
	}
	return sb.String()
}

func renderJSON(result comparator.Result) (string, error) {
	out := JSONOutput{
		MissingInRight: result.MissingInRight,
		MissingInLeft:  result.MissingInLeft,
		Mismatched:     result.Mismatched,
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b), nil
}

func renderMarkdown(result comparator.Result) string {
	var sb strings.Builder
	if len(result.MissingInRight) > 0 {
		sb.WriteString("### Missing in right\n")
		for _, k := range result.MissingInRight {
			fmt.Fprintf(&sb, "- `%s`\n", k)
		}
	}
	if len(result.MissingInLeft) > 0 {
		sb.WriteString("### Missing in left\n")
		for _, k := range result.MissingInLeft {
			fmt.Fprintf(&sb, "- `%s`\n", k)
		}
	}
	if len(result.Mismatched) > 0 {
		sb.WriteString("### Mismatched values\n")
		for k, v := range result.Mismatched {
			fmt.Fprintf(&sb, "- `%s`: `%s` → `%s`\n", k, v[0], v[1])
		}
	}
	return sb.String()
}
