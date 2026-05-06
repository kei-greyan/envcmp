package stats_test

import (
	"testing"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/stats"
)

func baseResult() comparator.Result {
	return comparator.Result{
		MissingInLeft:  []string{},
		MissingInRight: []string{},
		Mismatched:     map[string][2]string{},
		Matched:        map[string]string{},
	}
}

func TestCompute_EmptyResult(t *testing.T) {
	s := stats.Compute(baseResult())
	if s.TotalKeys != 0 || s.HasDiff() {
		t.Errorf("expected empty summary, got %+v", s)
	}
}

func TestCompute_OnlyMatched(t *testing.T) {
	r := baseResult()
	r.Matched = map[string]string{"HOST": "localhost", "PORT": "8080"}
	s := stats.Compute(r)
	if s.TotalKeys != 2 || s.Matched != 2 || s.HasDiff() {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestCompute_MissingInRight(t *testing.T) {
	r := baseResult()
	r.MissingInRight = []string{"SECRET", "TOKEN"}
	s := stats.Compute(r)
	if s.MissingInRight != 2 || s.TotalKeys != 2 {
		t.Errorf("unexpected summary: %+v", s)
	}
	if !s.HasDiff() {
		t.Error("expected HasDiff to be true")
	}
}

func TestCompute_MissingInLeft(t *testing.T) {
	r := baseResult()
	r.MissingInLeft = []string{"DB_URL"}
	s := stats.Compute(r)
	if s.MissingInLeft != 1 || !s.HasDiff() {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestCompute_Mismatched(t *testing.T) {
	r := baseResult()
	r.Mismatched = map[string][2]string{
		"LOG_LEVEL": {"debug", "info"},
		"TIMEOUT":   {"30", "60"},
	}
	s := stats.Compute(r)
	if s.Mismatched != 2 || s.TotalKeys != 2 || !s.HasDiff() {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestCompute_TotalKeys_NoDuplicates(t *testing.T) {
	r := baseResult()
	r.Matched = map[string]string{"A": "1", "B": "2"}
	r.MissingInRight = []string{"C"}
	r.MissingInLeft = []string{"D"}
	r.Mismatched = map[string][2]string{"E": {"x", "y"}}
	s := stats.Compute(r)
	if s.TotalKeys != 5 {
		t.Errorf("expected 5 total keys, got %d", s.TotalKeys)
	}
}
