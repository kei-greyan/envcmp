package envmerge

import (
	"strings"
	"testing"
)

func TestFormat_NoEntries(t *testing.T) {
	r := Result{}
	out := Format(r)
	if !strings.Contains(out, "no keys merged") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestFormat_NoConflicts(t *testing.T) {
	r := Result{
		Entries: []Entry{
			{Key: "FOO", Value: "1", Source: 0, Conflict: false},
		},
	}
	out := Format(r)
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %q", out)
	}
	if strings.Contains(out, "conflict") {
		t.Errorf("unexpected conflict marker: %q", out)
	}
}

func TestFormat_WithConflict(t *testing.T) {
	r := Result{
		Entries: []Entry{
			{Key: "KEY", Value: "override", Source: 1, Conflict: true},
		},
	}
	out := Format(r)
	if !strings.Contains(out, "~") {
		t.Errorf("expected conflict marker '~', got: %q", out)
	}
	if !strings.Contains(out, "conflict") {
		t.Errorf("expected word 'conflict', got: %q", out)
	}
}

func TestSummary_NoConflicts(t *testing.T) {
	r := Result{
		Entries: []Entry{
			{Key: "A", Conflict: false},
			{Key: "B", Conflict: false},
		},
	}
	s := Summary(r)
	if !strings.Contains(s, "no conflicts") {
		t.Errorf("expected no conflicts message, got: %q", s)
	}
}

func TestSummary_WithConflicts(t *testing.T) {
	r := Result{
		Entries: []Entry{
			{Key: "A", Conflict: true},
			{Key: "B", Conflict: false},
		},
	}
	s := Summary(r)
	if !strings.Contains(s, "1 conflict") {
		t.Errorf("expected conflict count, got: %q", s)
	}
}
