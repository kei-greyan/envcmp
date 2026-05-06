package schema_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/schema"
)

func TestFormat_NoViolations(t *testing.T) {
	out := schema.Format(nil)
	if !strings.Contains(out, "no violations") {
		t.Errorf("expected 'no violations', got %q", out)
	}
}

func TestFormat_SingleViolation(t *testing.T) {
	v := []schema.Violation{{Key: "PORT", Message: "expected int"}}
	out := schema.Format(v)
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got %q", out)
	}
	if !strings.Contains(out, "expected int") {
		t.Errorf("expected message in output, got %q", out)
	}
	if !strings.Contains(out, "1 violation") {
		t.Errorf("expected count in output, got %q", out)
	}
}

func TestFormat_MultipleViolations_SortedByKey(t *testing.T) {
	v := []schema.Violation{
		{Key: "Z_KEY", Message: "missing"},
		{Key: "A_KEY", Message: "bad type"},
	}
	out := schema.Format(v)
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA == -1 || idxZ == -1 {
		t.Fatal("expected both keys in output")
	}
	if idxA > idxZ {
		t.Error("expected A_KEY before Z_KEY (sorted)")
	}
}

func TestSummary_NoViolations(t *testing.T) {
	out := schema.Summary(nil)
	if out != "OK" {
		t.Errorf("expected OK, got %q", out)
	}
}

func TestSummary_WithViolations(t *testing.T) {
	v := []schema.Violation{
		{Key: "A", Message: "x"},
		{Key: "B", Message: "y"},
	}
	out := schema.Summary(v)
	if !strings.Contains(out, "2") {
		t.Errorf("expected count 2 in summary, got %q", out)
	}
}
