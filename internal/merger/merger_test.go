package merger_test

import (
	"testing"

	"github.com/user/envcmp/internal/merger"
)

func TestMerge_NoConflicts(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"C": "3", "D": "4"}

	res, err := merger.Merge(left, right, merger.PreferLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Merged) != 4 {
		t.Errorf("expected 4 keys, got %d", len(res.Merged))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_PreferLeft_OnConflict(t *testing.T) {
	left := map[string]string{"KEY": "left-value"}
	right := map[string]string{"KEY": "right-value"}

	res, err := merger.Merge(left, right, merger.PreferLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "left-value" {
		t.Errorf("expected left-value, got %q", res.Merged["KEY"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "KEY" {
		t.Errorf("expected conflict on KEY, got %v", res.Conflicts)
	}
}

func TestMerge_PreferRight_OnConflict(t *testing.T) {
	left := map[string]string{"KEY": "left-value"}
	right := map[string]string{"KEY": "right-value"}

	res, err := merger.Merge(left, right, merger.PreferRight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "right-value" {
		t.Errorf("expected right-value, got %q", res.Merged["KEY"])
	}
}

func TestMerge_ErrorOnConflict_ReturnsError(t *testing.T) {
	left := map[string]string{"KEY": "a"}
	right := map[string]string{"KEY": "b"}

	_, err := merger.Merge(left, right, merger.ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_SameValue_NoConflict(t *testing.T) {
	left := map[string]string{"KEY": "same"}
	right := map[string]string{"KEY": "same"}

	res, err := merger.Merge(left, right, merger.ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts for equal values, got %v", res.Conflicts)
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	res, err := merger.Merge(map[string]string{}, map[string]string{}, merger.PreferLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Merged) != 0 {
		t.Errorf("expected empty merged map, got %d keys", len(res.Merged))
	}
}
