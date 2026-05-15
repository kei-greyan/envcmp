package envscope_test

import (
	"testing"

	"github.com/user/envcmp/internal/envscope"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.local",
		"DB_PORT":     "5432",
		"LOG_LEVEL":   "info",
		"APP_TIMEOUT": "30s",
	}
}

func TestExtract_ByPrefix_ReturnsMatchingKeys(t *testing.T) {
	r, err := envscope.Extract(baseEnv(), envscope.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Matched != 3 {
		t.Errorf("expected 3 matched, got %d", r.Matched)
	}
	if _, ok := r.Env["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := r.Env["DB_HOST"]; ok {
		t.Error("did not expect DB_HOST in result")
	}
}

func TestExtract_StripPrefix_RemovesPrefixFromKeys(t *testing.T) {
	r, err := envscope.Extract(baseEnv(), envscope.Options{Prefix: "APP_", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Env["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_ prefix")
	}
	if _, ok := r.Env["APP_HOST"]; ok {
		t.Error("did not expect APP_HOST after stripping prefix")
	}
}

func TestExtract_BySuffix_ReturnsMatchingKeys(t *testing.T) {
	r, err := envscope.Extract(baseEnv(), envscope.Options{Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Matched != 2 {
		t.Errorf("expected 2 matched, got %d", r.Matched)
	}
}

func TestExtract_PrefixAndSuffix_BothMustMatch(t *testing.T) {
	r, err := envscope.Extract(baseEnv(), envscope.Options{Prefix: "APP_", Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Matched != 1 {
		t.Errorf("expected 1 matched, got %d", r.Matched)
	}
	if _, ok := r.Env["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestExtract_NoOptions_ReturnsError(t *testing.T) {
	_, err := envscope.Extract(baseEnv(), envscope.Options{})
	if err == nil {
		t.Error("expected error when no prefix or suffix set")
	}
}

func TestExtract_ExcludedCount_IsCorrect(t *testing.T) {
	r, err := envscope.Extract(baseEnv(), envscope.Options{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Excluded != 4 {
		t.Errorf("expected 4 excluded, got %d", r.Excluded)
	}
}

func TestKeys_ReturnsSortedKeys(t *testing.T) {
	r, _ := envscope.Extract(baseEnv(), envscope.Options{Prefix: "APP_"})
	keys := envscope.Keys(r)
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted at index %d: %s < %s", i, keys[i], keys[i-1])
		}
	}
}
