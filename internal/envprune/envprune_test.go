package envprune_test

import (
	"testing"

	"github.com/user/envcmp/internal/envprune"
)

var baseEnv = map[string]string{
	"APP_NAME":    "myapp",
	"APP_SECRET":  "s3cr3t",
	"DEBUG":       "true",
	"EMPTY_KEY":   "",
	"OLD_FEATURE": "1",
	"LOG_LEVEL":   "info",
}

func TestPrune_RemoveEmpty(t *testing.T) {
	r := envprune.Prune(baseEnv, envprune.Options{RemoveEmpty: true})
	if _, ok := r.Pruned["EMPTY_KEY"]; !ok {
		t.Error("expected EMPTY_KEY to be pruned")
	}
	if _, ok := r.Kept["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to be kept")
	}
}

func TestPrune_RemoveByPrefix(t *testing.T) {
	r := envprune.Prune(baseEnv, envprune.Options{RemovePrefixes: []string{"APP_"}})
	if _, ok := r.Pruned["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to be pruned")
	}
	if _, ok := r.Pruned["APP_SECRET"]; !ok {
		t.Error("expected APP_SECRET to be pruned")
	}
	if _, ok := r.Kept["DEBUG"]; !ok {
		t.Error("expected DEBUG to be kept")
	}
}

func TestPrune_RemoveBySuffix(t *testing.T) {
	r := envprune.Prune(baseEnv, envprune.Options{RemoveSuffixes: []string{"_LEVEL"}})
	if _, ok := r.Pruned["LOG_LEVEL"]; !ok {
		t.Error("expected LOG_LEVEL to be pruned")
	}
	if _, ok := r.Kept["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to be kept")
	}
}

func TestPrune_RemoveExactKeys(t *testing.T) {
	r := envprune.Prune(baseEnv, envprune.Options{RemoveKeys: []string{"DEBUG", "OLD_FEATURE"}})
	if _, ok := r.Pruned["DEBUG"]; !ok {
		t.Error("expected DEBUG to be pruned")
	}
	if _, ok := r.Pruned["OLD_FEATURE"]; !ok {
		t.Error("expected OLD_FEATURE to be pruned")
	}
	if len(r.Kept)+len(r.Pruned) != len(baseEnv) {
		t.Errorf("total keys mismatch: got %d, want %d", len(r.Kept)+len(r.Pruned), len(baseEnv))
	}
}

func TestPrune_NoOptions_KeepsAll(t *testing.T) {
	r := envprune.Prune(baseEnv, envprune.Options{})
	if len(r.Pruned) != 0 {
		t.Errorf("expected no pruned keys, got %d", len(r.Pruned))
	}
	if len(r.Kept) != len(baseEnv) {
		t.Errorf("expected all keys kept, got %d", len(r.Kept))
	}
}

func TestPrune_DoesNotMutateOriginal(t *testing.T) {
	origLen := len(baseEnv)
	envprune.Prune(baseEnv, envprune.Options{RemoveEmpty: true, RemovePrefixes: []string{"APP_"}})
	if len(baseEnv) != origLen {
		t.Error("original env map was mutated")
	}
}

func TestSummary_Format(t *testing.T) {
	r := envprune.Prune(baseEnv, envprune.Options{RemoveEmpty: true})
	s := envprune.Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
	f := envprune.Format(r)
	if f == "" {
		t.Error("expected non-empty format output")
	}
}
