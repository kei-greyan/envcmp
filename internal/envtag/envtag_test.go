package envtag_test

import (
	"testing"

	"github.com/user/envcmp/internal/envtag"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"AUTH_SECRET": "s3cr3t",
	"JWT_TTL":     "3600",
	"APP_NAME":    "envcmp",
	"LOG_LEVEL":   "info",
}

func TestApply_GroupsByPrefix(t *testing.T) {
	opts := envtag.Options{
		Tags: map[string][]string{
			"database": {"DB_"},
			"auth":     {"AUTH_", "JWT_"},
		},
	}
	r := envtag.Apply(baseEnv, opts)

	if len(r.Tagged["database"]) != 2 {
		t.Errorf("expected 2 database keys, got %d", len(r.Tagged["database"]))
	}
	if len(r.Tagged["auth"]) != 2 {
		t.Errorf("expected 2 auth keys, got %d", len(r.Tagged["auth"]))
	}
	if len(r.Untagged) != 2 {
		t.Errorf("expected 2 untagged keys, got %d", len(r.Untagged))
	}
}

func TestApply_ExactMatch(t *testing.T) {
	opts := envtag.Options{
		Tags: map[string][]string{
			"specific": {"DB_HOST", "LOG_LEVEL"},
		},
		ExactMatch: true,
	}
	r := envtag.Apply(baseEnv, opts)

	if len(r.Tagged["specific"]) != 2 {
		t.Errorf("expected 2 specific keys, got %d", len(r.Tagged["specific"]))
	}
	if len(r.Untagged) != 4 {
		t.Errorf("expected 4 untagged, got %d", len(r.Untagged))
	}
}

func TestApply_NoTags_AllUntagged(t *testing.T) {
	r := envtag.Apply(baseEnv, envtag.Options{})

	if len(r.Tagged) != 0 {
		t.Errorf("expected no tagged groups, got %d", len(r.Tagged))
	}
	if len(r.Untagged) != len(baseEnv) {
		t.Errorf("expected all keys untagged, got %d", len(r.Untagged))
	}
}

func TestApply_EmptyEnv_ReturnsEmptyResult(t *testing.T) {
	opts := envtag.Options{
		Tags: map[string][]string{"db": {"DB_"}},
	}
	r := envtag.Apply(map[string]string{}, opts)

	if len(r.Tagged["db"]) != 0 {
		t.Errorf("expected 0 tagged keys")
	}
	if len(r.Untagged) != 0 {
		t.Errorf("expected 0 untagged keys")
	}
}

func TestTagNames_ReturnsSorted(t *testing.T) {
	r := envtag.Result{
		Tagged: map[string][]string{
			"zebra": {"Z_KEY"},
			"alpha": {"A_KEY"},
			"mango": {"M_KEY"},
		},
	}
	names := envtag.TagNames(r)
	expected := []string{"alpha", "mango", "zebra"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("pos %d: expected %s got %s", i, expected[i], n)
		}
	}
}

func TestSummary_Output(t *testing.T) {
	opts := envtag.Options{
		Tags: map[string][]string{"db": {"DB_"}},
	}
	r := envtag.Apply(baseEnv, opts)
	s := envtag.Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
