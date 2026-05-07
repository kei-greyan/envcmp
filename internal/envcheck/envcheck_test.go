package envcheck_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envcheck"
)

var reference = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"SECRET_KEY":  "changeme",
	"FEATURE_FLAG": "false",
}

func TestCheck_AllPresent(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "prod.db",
		"DB_PORT":     "5432",
		"SECRET_KEY":  "s3cr3t",
		"FEATURE_FLAG": "true",
	}
	r := envcheck.Check(reference, env, envcheck.Options{})
	if !r.IsClean() {
		t.Fatalf("expected clean result, got missing=%v empty=%v", r.Missing, r.Empty)
	}
	if len(r.Present) != 4 {
		t.Errorf("expected 4 present, got %d", len(r.Present))
	}
}

func TestCheck_MissingKey(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "prod.db",
		"DB_PORT": "5432",
	}
	r := envcheck.Check(reference, env, envcheck.Options{})
	if r.IsClean() {
		t.Fatal("expected dirty result")
	}
	if len(r.Missing) != 2 {
		t.Errorf("expected 2 missing, got %d", len(r.Missing))
	}
}

func TestCheck_EmptyValue_FailOnEmpty(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "prod.db",
		"DB_PORT":     "5432",
		"SECRET_KEY":  "",
		"FEATURE_FLAG": "true",
	}
	r := envcheck.Check(reference, env, envcheck.Options{FailOnEmpty: true})
	if r.IsClean() {
		t.Fatal("expected dirty result due to empty key")
	}
	if len(r.Empty) != 1 || r.Empty[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY in empty list, got %v", r.Empty)
	}
}

func TestCheck_EmptyValue_NoFailOnEmpty(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "prod.db",
		"DB_PORT":     "5432",
		"SECRET_KEY":  "",
		"FEATURE_FLAG": "true",
	}
	r := envcheck.Check(reference, env, envcheck.Options{FailOnEmpty: false})
	if !r.IsClean() {
		t.Fatalf("expected clean result, got %+v", r)
	}
}

func TestFormat_CleanResult(t *testing.T) {
	r := envcheck.Result{Present: []string{"A", "B"}}
	out := envcheck.Format(r)
	if !strings.HasPrefix(out, "OK:") {
		t.Errorf("expected OK prefix, got %q", out)
	}
}

func TestFormat_DirtyResult(t *testing.T) {
	r := envcheck.Result{Missing: []string{"DB_HOST"}, Empty: []string{"SECRET_KEY"}}
	out := envcheck.Format(r)
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING in output, got %q", out)
	}
	if !strings.Contains(out, "EMPTY") {
		t.Errorf("expected EMPTY in output, got %q", out)
	}
}

func TestSummary(t *testing.T) {
	r := envcheck.Result{
		Present: []string{"A", "B"},
		Missing: []string{"C"},
		Empty:   []string{"D"},
	}
	s := envcheck.Summary(r)
	if !strings.Contains(s, "2 present") || !strings.Contains(s, "1 missing") || !strings.Contains(s, "1 empty") {
		t.Errorf("unexpected summary: %q", s)
	}
}
