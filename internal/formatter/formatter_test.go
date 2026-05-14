package formatter_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/formatter"
)

var baseResult = comparator.Result{
	MissingInRight: []string{"FOO"},
	MissingInLeft:  []string{"BAR"},
	Mismatched:     map[string][2]string{"BAZ": {"old", "new"}},
}

func TestRender_TextFormat(t *testing.T) {
	out, err := formatter.Render(baseResult, formatter.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[MISSING RIGHT] FOO") {
		t.Errorf("expected MISSING RIGHT FOO, got: %s", out)
	}
	if !strings.Contains(out, "[MISSING LEFT]  BAR") {
		t.Errorf("expected MISSING LEFT BAR, got: %s", out)
	}
	if !strings.Contains(out, "[MISMATCH]      BAZ") {
		t.Errorf("expected MISMATCH BAZ, got: %s", out)
	}
}

func TestRender_JSONFormat(t *testing.T) {
	out, err := formatter.Render(baseResult, formatter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed formatter.JSONOutput
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(parsed.MissingInRight) != 1 || parsed.MissingInRight[0] != "FOO" {
		t.Errorf("unexpected MissingInRight: %v", parsed.MissingInRight)
	}
	if len(parsed.MissingInLeft) != 1 || parsed.MissingInLeft[0] != "BAR" {
		t.Errorf("unexpected MissingInLeft: %v", parsed.MissingInLeft)
	}
	if v, ok := parsed.Mismatched["BAZ"]; !ok || v[0] != "old" || v[1] != "new" {
		t.Errorf("unexpected Mismatched: %v", parsed.Mismatched)
	}
}

func TestRender_MarkdownFormat(t *testing.T) {
	out, err := formatter.Render(baseResult, formatter.FormatMarkdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "### Missing in right") {
		t.Errorf("expected markdown header, got: %s", out)
	}
	if !strings.Contains(out, "`FOO`") {
		t.Errorf("expected FOO in markdown, got: %s", out)
	}
}

func TestRender_UnknownFormat(t *testing.T) {
	_, err := formatter.Render(baseResult, formatter.Format("xml"))
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestRender_EmptyResult_JSON(t *testing.T) {
	out, err := formatter.Render(comparator.Result{}, formatter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "{") {
		t.Errorf("expected valid JSON object, got: %s", out)
	}
}

func TestRender_EmptyResult_Text(t *testing.T) {
	out, err := formatter.Render(comparator.Result{}, formatter.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// An empty result should produce no diff lines.
	if strings.Contains(out, "[MISSING") || strings.Contains(out, "[MISMATCH") {
		t.Errorf("expected empty text output for empty result, got: %s", out)
	}
}
