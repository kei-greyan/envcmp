package envaudit_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envaudit"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "envcmp",
	}
}

func TestAudit_NoChanges(t *testing.T) {
	before := baseEnv()
	after := baseEnv()
	r := envaudit.Audit(before, after)
	if !r.IsClean() {
		t.Fatalf("expected clean report, got %d entries", len(r.Entries))
	}
}

func TestAudit_AddedKey(t *testing.T) {
	before := baseEnv()
	after := baseEnv()
	after["NEW_KEY"] = "value"

	r := envaudit.Audit(before, after)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Kind != envaudit.Added || r.Entries[0].Key != "NEW_KEY" {
		t.Errorf("unexpected entry: %+v", r.Entries[0])
	}
}

func TestAudit_RemovedKey(t *testing.T) {
	before := baseEnv()
	after := baseEnv()
	delete(after, "DB_PORT")

	r := envaudit.Audit(before, after)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Kind != envaudit.Removed || r.Entries[0].Key != "DB_PORT" {
		t.Errorf("unexpected entry: %+v", r.Entries[0])
	}
	if r.Entries[0].OldValue != "5432" {
		t.Errorf("expected OldValue=5432, got %q", r.Entries[0].OldValue)
	}
}

func TestAudit_ModifiedKey(t *testing.T) {
	before := baseEnv()
	after := baseEnv()
	after["APP_ENV"] = "staging"

	r := envaudit.Audit(before, after)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Kind != envaudit.Modified || e.Key != "APP_ENV" {
		t.Errorf("unexpected entry: %+v", e)
	}
	if e.OldValue != "production" || e.NewValue != "staging" {
		t.Errorf("unexpected values: old=%q new=%q", e.OldValue, e.NewValue)
	}
}

func TestAudit_SortedByKey(t *testing.T) {
	before := map[string]string{"Z_KEY": "1", "A_KEY": "2"}
	after := map[string]string{"Z_KEY": "changed", "A_KEY": "changed"}

	r := envaudit.Audit(before, after)
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
	if r.Entries[0].Key != "A_KEY" || r.Entries[1].Key != "Z_KEY" {
		t.Errorf("entries not sorted: %v, %v", r.Entries[0].Key, r.Entries[1].Key)
	}
}

func TestFormat_CleanReport(t *testing.T) {
	r := envaudit.Report{}
	out := envaudit.Format(r)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' in output, got %q", out)
	}
}

func TestFormat_WithChanges(t *testing.T) {
	before := map[string]string{"OLD": "val"}
	after := map[string]string{"NEW": "val"}
	r := envaudit.Audit(before, after)
	out := envaudit.Format(r)
	if !strings.Contains(out, "NEW") || !strings.Contains(out, "OLD") {
		t.Errorf("expected keys in output, got %q", out)
	}
}

func TestSummary_NoChanges(t *testing.T) {
	r := envaudit.Report{}
	if envaudit.Summary(r) != "no changes" {
		t.Errorf("unexpected summary: %q", envaudit.Summary(r))
	}
}

func TestSummary_WithChanges(t *testing.T) {
	before := baseEnv()
	after := baseEnv()
	after["APP_ENV"] = "staging"
	after["EXTRA"] = "new"
	delete(after, "DB_PORT")

	r := envaudit.Audit(before, after)
	s := envaudit.Summary(r)
	if !strings.Contains(s, "1 added") || !strings.Contains(s, "1 removed") || !strings.Contains(s, "1 modified") {
		t.Errorf("unexpected summary: %q", s)
	}
}
