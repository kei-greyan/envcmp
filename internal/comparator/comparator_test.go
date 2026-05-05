package comparator

import (
	"testing"
)

func TestCompare_NoDiff(t *testing.T) {
	left := map[string]string{"KEY": "value", "FOO": "bar"}
	right := map[string]string{"KEY": "value", "FOO": "bar"}

	res := Compare(left, right)
	if res.HasDiff() {
		t.Errorf("expected no diff, got %+v", res)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"KEY": "value", "ONLY_LEFT": "x"}
	right := map[string]string{"KEY": "value"}

	res := Compare(left, right)
	if len(res.MissingInRight) != 1 || res.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("expected ONLY_LEFT missing in right, got %v", res.MissingInRight)
	}
	if len(res.MissingInLeft) != 0 {
		t.Errorf("expected no missing in left, got %v", res.MissingInLeft)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"KEY": "value"}
	right := map[string]string{"KEY": "value", "ONLY_RIGHT": "y"}

	res := Compare(left, right)
	if len(res.MissingInLeft) != 1 || res.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Errorf("expected ONLY_RIGHT missing in left, got %v", res.MissingInLeft)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := map[string]string{"KEY": "old", "SAME": "same"}
	right := map[string]string{"KEY": "new", "SAME": "same"}

	res := Compare(left, right)
	if len(res.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(res.Mismatched))
	}
	m := res.Mismatched[0]
	if m.Key != "KEY" || m.LeftValue != "old" || m.RightValue != "new" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	res := Compare(map[string]string{}, map[string]string{})
	if res.HasDiff() {
		t.Error("expected no diff for two empty maps")
	}
}

func TestHasDiff(t *testing.T) {
	res := Result{MissingInRight: []string{"A"}}
	if !res.HasDiff() {
		t.Error("expected HasDiff to return true")
	}
}
