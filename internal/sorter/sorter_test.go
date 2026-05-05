package sorter_test

import (
	"testing"

	"github.com/yourusername/envcmp/internal/comparator"
	"github.com/yourusername/envcmp/internal/sorter"
)

func baseResult() comparator.Result {
	return comparator.Result{
		MissingInRight: map[string]string{
			"ZEBRA": "z",
			"APPLE": "a",
			"MANGO": "m",
		},
		MissingInLeft: map[string]string{
			"BANANA": "b",
			"CHERRY": "c",
		},
		Mismatched: map[string]comparator.ValuePair{
			"PORT": {Left: "8080", Right: "9090"},
			"HOST": {Left: "localhost", Right: "prod.example.com"},
		},
	}
}

func TestSort_MissingInRight_IsSorted(t *testing.T) {
	sr := sorter.Sort(baseResult())
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	if len(sr.MissingInRight) != len(expected) {
		t.Fatalf("expected %d entries, got %d", len(expected), len(sr.MissingInRight))
	}
	for i, key := range expected {
		if sr.MissingInRight[i] != key {
			t.Errorf("index %d: expected %q, got %q", i, key, sr.MissingInRight[i])
		}
	}
}

func TestSort_MissingInLeft_IsSorted(t *testing.T) {
	sr := sorter.Sort(baseResult())
	expected := []string{"BANANA", "CHERRY"}
	for i, key := range expected {
		if sr.MissingInLeft[i] != key {
			t.Errorf("index %d: expected %q, got %q", i, key, sr.MissingInLeft[i])
		}
	}
}

func TestSort_Mismatched_IsSorted(t *testing.T) {
	sr := sorter.Sort(baseResult())
	if len(sr.Mismatched) != 2 {
		t.Fatalf("expected 2 mismatched entries, got %d", len(sr.Mismatched))
	}
	if sr.Mismatched[0].Key != "HOST" {
		t.Errorf("expected first key HOST, got %q", sr.Mismatched[0].Key)
	}
	if sr.Mismatched[1].Key != "PORT" {
		t.Errorf("expected second key PORT, got %q", sr.Mismatched[1].Key)
	}
}

func TestSort_MismatchedValues_Preserved(t *testing.T) {
	sr := sorter.Sort(baseResult())
	for _, entry := range sr.Mismatched {
		if entry.Key == "PORT" {
			if entry.Left != "8080" || entry.Right != "9090" {
				t.Errorf("PORT values not preserved: got left=%q right=%q", entry.Left, entry.Right)
			}
		}
	}
}

func TestSort_EmptyResult(t *testing.T) {
	sr := sorter.Sort(comparator.Result{})
	if len(sr.MissingInRight) != 0 || len(sr.MissingInLeft) != 0 || len(sr.Mismatched) != 0 {
		t.Error("expected all empty slices for empty result")
	}
}
