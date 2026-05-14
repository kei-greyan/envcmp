package envmerge

import (
	"strings"
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	r, err := Merge(StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["FOO"] != "1" || r.Env["BAR"] != "2" || r.Env["BAZ"] != "3" {
		t.Errorf("unexpected env: %v", r.Env)
	}
	for _, e := range r.Entries {
		if e.Conflict {
			t.Errorf("unexpected conflict on %s", e.Key)
		}
	}
}

func TestMerge_StrategyFirst_KeepsOriginal(t *testing.T) {
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "override"}

	r, err := Merge(StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["KEY"] != "original" {
		t.Errorf("expected 'original', got %q", r.Env["KEY"])
	}
	if !r.Entries[0].Conflict {
		t.Error("expected conflict flag on KEY")
	}
}

func TestMerge_StrategyLast_OverridesValue(t *testing.T) {
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "override"}

	r, err := Merge(StrategyLast, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["KEY"] != "override" {
		t.Errorf("expected 'override', got %q", r.Env["KEY"])
	}
}

func TestMerge_StrategyError_ReturnsError(t *testing.T) {
	a := map[string]string{"KEY": "a"}
	b := map[string]string{"KEY": "b"}

	_, err := Merge(StrategyError, a, b)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "KEY") {
		t.Errorf("error should mention key, got: %v", err)
	}
}

func TestMerge_SameValue_NoConflict(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}

	r, err := Merge(StrategyError, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Entries[0].Conflict {
		t.Error("identical values should not be flagged as conflict")
	}
}

func TestMerge_EntriesSortedByKey(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "2", "M": "3"}

	r, err := Merge(StrategyFirst, a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("entries not sorted: %v", keys)
	}
}
