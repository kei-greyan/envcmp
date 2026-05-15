package envdrift_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envdrift"
)

func baseBaseline() map[string]string {
	return map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "db.prod",
		"SECRET":  "abc123",
	}
}

func TestDetect_NoDrift(t *testing.T) {
	b := baseBaseline()
	c := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "db.prod",
		"SECRET":  "abc123",
	}
	r := envdrift.Detect(b, c)
	if r.HasDrift() {
		t.Fatalf("expected no drift, got %d entries", len(r.Entries))
	}
}

func TestDetect_AddedKey(t *testing.T) {
	b := baseBaseline()
	c := baseBaseline()
	c["NEW_KEY"] = "new"

	r := envdrift.Detect(b, c)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Kind != envdrift.DriftAdded || r.Entries[0].Key != "NEW_KEY" {
		t.Errorf("unexpected entry: %+v", r.Entries[0])
	}
}

func TestDetect_RemovedKey(t *testing.T) {
	b := baseBaseline()
	c := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "db.prod",
	}

	r := envdrift.Detect(b, c)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Kind != envdrift.DriftRemoved || r.Entries[0].Key != "SECRET" {
		t.Errorf("unexpected entry: %+v", r.Entries[0])
	}
}

func TestDetect_ChangedKey(t *testing.T) {
	b := baseBaseline()
	c := baseBaseline()
	c["DB_HOST"] = "db.staging"

	r := envdrift.Detect(b, c)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Kind != envdrift.DriftChanged {
		t.Errorf("expected changed, got %s", e.Kind)
	}
	if e.Baseline != "db.prod" || e.Current != "db.staging" {
		t.Errorf("unexpected values: baseline=%q current=%q", e.Baseline, e.Current)
	}
}

func TestDetect_SortedByKey(t *testing.T) {
	b := map[string]string{"Z_KEY": "1", "A_KEY": "2"}
	c := map[string]string{"Z_KEY": "changed", "A_KEY": "changed"}

	r := envdrift.Detect(b, c)
	if r.Entries[0].Key > r.Entries[1].Key {
		t.Errorf("entries not sorted: %s > %s", r.Entries[0].Key, r.Entries[1].Key)
	}
}

func TestSummary_NoDrift(t *testing.T) {
	r := envdrift.Result{}
	if got := envdrift.Summary(r); got != "no drift detected" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_WithDrift(t *testing.T) {
	r := envdrift.Result{
		Entries: []envdrift.Entry{
			{Key: "A", Kind: envdrift.DriftAdded},
			{Key: "B", Kind: envdrift.DriftRemoved},
			{Key: "C", Kind: envdrift.DriftChanged},
		},
	}
	s := envdrift.Summary(r)
	if !strings.Contains(s, "1 added") || !strings.Contains(s, "1 removed") || !strings.Contains(s, "1 changed") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestFormat_NoDrift(t *testing.T) {
	out := envdrift.Format(envdrift.Result{})
	if !strings.Contains(out, "No drift") {
		t.Errorf("expected no-drift message, got: %q", out)
	}
}

func TestFormat_ShowsKindSymbols(t *testing.T) {
	r := envdrift.Result{
		Entries: []envdrift.Entry{
			{Key: "ADDED_KEY", Kind: envdrift.DriftAdded, Current: "v"},
			{Key: "REMOVED_KEY", Kind: envdrift.DriftRemoved, Baseline: "v"},
			{Key: "CHANGED_KEY", Kind: envdrift.DriftChanged, Baseline: "old", Current: "new"},
		},
	}
	out := envdrift.Format(r)
	if !strings.Contains(out, "+") {
		t.Error("expected '+' symbol for added key")
	}
	if !strings.Contains(out, "-") {
		t.Error("expected '-' symbol for removed key")
	}
	if !strings.Contains(out, "~") {
		t.Error("expected '~' symbol for changed key")
	}
}
