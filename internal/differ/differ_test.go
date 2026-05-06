package differ_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/differ"
)

func TestDiffValues_NoDiff(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	got := differ.DiffValues(left, right)
	if len(got) != 0 {
		t.Fatalf("expected no diffs, got %d", len(got))
	}
}

func TestDiffValues_SingleChange(t *testing.T) {
	left := map[string]string{"A": "old", "B": "same"}
	right := map[string]string{"A": "new", "B": "same"}
	got := differ.DiffValues(left, right)
	if len(got) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(got))
	}
	if got[0].Key != "A" || got[0].Left != "old" || got[0].Right != "new" {
		t.Errorf("unexpected diff: %+v", got[0])
	}
}

func TestDiffValues_IgnoresMissingKeys(t *testing.T) {
	left := map[string]string{"A": "1", "ONLY_LEFT": "x"}
	right := map[string]string{"A": "1", "ONLY_RIGHT": "y"}
	got := differ.DiffValues(left, right)
	if len(got) != 0 {
		t.Fatalf("expected no diffs for shared keys, got %d", len(got))
	}
}

func TestDiffValues_SortedByKey(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "1", "M": "1"}
	right := map[string]string{"Z": "2", "A": "2", "M": "2"}
	got := differ.DiffValues(left, right)
	if len(got) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(got))
	}
	if got[0].Key != "A" || got[1].Key != "M" || got[2].Key != "Z" {
		t.Errorf("diffs not sorted: %v %v %v", got[0].Key, got[1].Key, got[2].Key)
	}
}

func TestDiffValues_EmptyMaps(t *testing.T) {
	got := differ.DiffValues(map[string]string{}, map[string]string{})
	if len(got) != 0 {
		t.Fatalf("expected no diffs for empty maps, got %d", len(got))
	}
}

func TestValueDiff_Describe(t *testing.T) {
	d := differ.ValueDiff{Key: "FOO", Left: "bar", Right: "baz"}
	got := d.Describe()
	if !strings.Contains(got, "FOO") || !strings.Contains(got, "bar") || !strings.Contains(got, "baz") {
		t.Errorf("unexpected description: %s", got)
	}
}

func TestValueDiff_IsEmpty(t *testing.T) {
	if !(differ.ValueDiff{Key: "K", Left: "v", Right: "v"}.IsEmpty()) {
		t.Error("expected IsEmpty true for equal values")
	}
	if (differ.ValueDiff{Key: "K", Left: "a", Right: "b"}.IsEmpty()) {
		t.Error("expected IsEmpty false for different values")
	}
}

func TestSummary_NoDiffs(t *testing.T) {
	got := differ.Summary(nil)
	if got != "no value differences" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestSummary_WithDiffs(t *testing.T) {
	diffs := []differ.ValueDiff{
		{Key: "A", Left: "x", Right: "y"},
		{Key: "B", Left: "1", Right: "2"},
	}
	got := differ.Summary(diffs)
	if !strings.Contains(got, "A") || !strings.Contains(got, "B") {
		t.Errorf("summary missing keys: %s", got)
	}
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
