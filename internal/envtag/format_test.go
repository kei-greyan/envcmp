package envtag_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envtag"
)

func TestFormat_ContainsTagHeaders(t *testing.T) {
	opts := envtag.Options{
		Tags: map[string][]string{
			"database": {"DB_"},
		},
	}
	r := envtag.Apply(baseEnv, opts)
	out := envtag.Format(r)

	if !strings.Contains(out, "[database]") {
		t.Error("expected [database] header in output")
	}
	if !strings.Contains(out, "[untagged]") {
		t.Error("expected [untagged] header in output")
	}
}

func TestFormat_ContainsKeyNames(t *testing.T) {
	opts := envtag.Options{
		Tags: map[string][]string{
			"database": {"DB_"},
		},
	}
	r := envtag.Apply(baseEnv, opts)
	out := envtag.Format(r)

	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "DB_PORT") {
		t.Error("expected DB_PORT in output")
	}
}

func TestFormat_NoUntagged_OmitsUntaggedHeader(t *testing.T) {
	r := envtag.Result{
		Tagged:   map[string][]string{"all": {"DB_HOST"}},
		Untagged: []string{},
	}
	out := envtag.Format(r)
	if strings.Contains(out, "[untagged]") {
		t.Error("should not contain [untagged] when list is empty")
	}
}

func TestSummary_CleanCounts(t *testing.T) {
	r := envtag.Result{
		Tagged:   map[string][]string{"db": {"DB_HOST", "DB_PORT"}, "auth": {"AUTH_SECRET"}},
		Untagged: []string{"APP_NAME"},
	}
	s := envtag.Summary(r)
	if !strings.Contains(s, "2 tag(s)") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "3 tagged") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "1 untagged") {
		t.Errorf("unexpected summary: %s", s)
	}
}
