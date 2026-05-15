package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envcmp/internal/comparator"
)

func TestReport_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	Report(&buf, comparator.Result{}, "left.env", "right.env")
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestReport_MissingInRight(t *testing.T) {
	res := comparator.Result{MissingInRight: []string{"SECRET"}}
	var buf bytes.Buffer
	Report(&buf, res, "dev.env", "prod.env")
	out := buf.String()
	if !strings.Contains(out, "dev.env") {
		t.Error("expected left name in output")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected key SECRET in output")
	}
	if !strings.Contains(out, "- SECRET") {
		t.Error("expected missing marker for SECRET")
	}
}

func TestReport_MissingInLeft(t *testing.T) {
	res := comparator.Result{MissingInLeft: []string{"NEW_KEY"}}
	var buf bytes.Buffer
	Report(&buf, res, "dev.env", "prod.env")
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected addition marker, got: %s", out)
	}
}

func TestReport_Mismatched(t *testing.T) {
	res := comparator.Result{
		Mismatched: []comparator.MismatchedKey{
			{Key: "DB_HOST", LeftValue: "localhost", RightValue: "db.prod.example.com"},
		},
	}
	var buf bytes.Buffer
	Report(&buf, res, "dev.env", "prod.env")
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "localhost") {
		t.Error("expected left value in output")
	}
	if !strings.Contains(out, "db.prod.example.com") {
		t.Error("expected right value in output")
	}
}

func TestReport_MultipleDiffs(t *testing.T) {
	res := comparator.Result{
		MissingInRight: []string{"OLD_KEY"},
		MissingInLeft:  []string{"NEW_KEY"},
		Mismatched: []comparator.MismatchedKey{
			{Key: "PORT", LeftValue: "3000", RightValue: "8080"},
		},
	}
	var buf bytes.Buffer
	Report(&buf, res, "dev.env", "prod.env")
	out := buf.String()
	if strings.Contains(out, "No differences") {
		t.Error("expected differences to be reported, got no-diff message")
	}
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected missing marker for OLD_KEY, got: %s", out)
	}
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected addition marker for NEW_KEY, got: %s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
}
