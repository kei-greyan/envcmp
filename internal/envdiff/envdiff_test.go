package envdiff_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envdiff"
)

func TestDiff_NoChanges(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1", "B": "2"}
	r := envdiff.Diff(old, new)
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
	for _, e := range r.Entries {
		if e.Kind != envdiff.Same {
			t.Errorf("expected Same for key %s, got %s", e.Key, e.Kind)
		}
	}
}

func TestDiff_AddedKey(t *testing.T) {
	old := map[string]string{"A": "1"}
	new := map[string]string{"A": "1", "B": "2"}
	r := envdiff.Diff(old, new)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	found := findEntry(r, "B")
	if found == nil {
		t.Fatal("expected entry for key B")
	}
	if found.Kind != envdiff.Added {
		t.Errorf("expected Added, got %s", found.Kind)
	}
	if found.NewValue != "2" {
		t.Errorf("expected NewValue=2, got %s", found.NewValue)
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1"}
	r := envdiff.Diff(old, new)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	found := findEntry(r, "B")
	if found == nil {
		t.Fatal("expected entry for key B")
	}
	if found.Kind != envdiff.Removed {
		t.Errorf("expected Removed, got %s", found.Kind)
	}
	if found.OldValue != "2" {
		t.Errorf("expected OldValue=2, got %s", found.OldValue)
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	old := map[string]string{"A": "old"}
	new := map[string]string{"A": "new"}
	r := envdiff.Diff(old, new)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	found := findEntry(r, "A")
	if found == nil {
		t.Fatal("expected entry for key A")
	}
	if found.Kind != envdiff.Changed {
		t.Errorf("expected Changed, got %s", found.Kind)
	}
	if found.OldValue != "old" || found.NewValue != "new" {
		t.Errorf("unexpected values: old=%s new=%s", found.OldValue, found.NewValue)
	}
}

func TestDiff_SortedByKey(t *testing.T) {
	old := map[string]string{"Z": "1", "A": "2", "M": "3"}
	new := map[string]string{"Z": "1", "A": "2", "M": "3"}
	r := envdiff.Diff(old, new)
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}

func TestFormat_NoChanges(t *testing.T) {
	r := envdiff.Diff(map[string]string{"X": "1"}, map[string]string{"X": "1"})
	out := envdiff.Format(r)
	if out != "no changes detected" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_ShowsAddedRemovedChanged(t *testing.T) {
	old := map[string]string{"A": "old", "B": "gone"}
	new := map[string]string{"A": "new", "C": "fresh"}
	r := envdiff.Diff(old, new)
	out := envdiff.Format(r)
	if !strings.Contains(out, "~ A") {
		t.Errorf("expected changed marker for A, got: %s", out)
	}
	if !strings.Contains(out, "- B") {
		t.Errorf("expected removed marker for B, got: %s", out)
	}
	if !strings.Contains(out, "+ C") {
		t.Errorf("expected added marker for C, got: %s", out)
	}
}

func findEntry(r envdiff.Result, key string) *envdiff.Entry {
	for i := range r.Entries {
		if r.Entries[i].Key == key {
			return &r.Entries[i]
		}
	}
	return nil
}
