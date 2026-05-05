package filter_test

import (
	"testing"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/filter"
)

func baseResult() comparator.Result {
	return comparator.Result{
		MissingInRight: []string{"ONLY_LEFT"},
		MissingInLeft:  []string{"ONLY_RIGHT"},
		Mismatched: []comparator.Diff{
			{Key: "DB_HOST", LeftVal: "localhost", RightVal: "prod.db"},
		},
	}
}

func TestApply_NoFilter_ReturnsAll(t *testing.T) {
	out := filter.Apply(baseResult(), filter.Options{})
	if len(out.MissingInRight) != 1 || len(out.MissingInLeft) != 1 || len(out.Mismatched) != 1 {
		t.Errorf("expected all entries, got MissingInRight=%d MissingInLeft=%d Mismatched=%d",
			len(out.MissingInRight), len(out.MissingInLeft), len(out.Mismatched))
	}
}

func TestApply_OnlyMissing_ExcludesMismatched(t *testing.T) {
	out := filter.Apply(baseResult(), filter.Options{OnlyMissing: true})
	if len(out.Mismatched) != 0 {
		t.Errorf("expected no mismatched entries, got %d", len(out.Mismatched))
	}
	if len(out.MissingInRight) != 1 || len(out.MissingInLeft) != 1 {
		t.Error("expected missing entries to be present")
	}
}

func TestApply_OnlyMismatched_ExcludesMissing(t *testing.T) {
	out := filter.Apply(baseResult(), filter.Options{OnlyMismatched: true})
	if len(out.MissingInRight) != 0 || len(out.MissingInLeft) != 0 {
		t.Errorf("expected no missing entries, got MissingInRight=%d MissingInLeft=%d",
			len(out.MissingInRight), len(out.MissingInLeft))
	}
	if len(out.Mismatched) != 1 {
		t.Error("expected mismatched entries to be present")
	}
}

func TestApply_KeyFilter_MatchingKey(t *testing.T) {
	out := filter.Apply(baseResult(), filter.Options{Keys: []string{"DB_HOST"}})
	if len(out.Mismatched) != 1 {
		t.Errorf("expected 1 mismatched entry, got %d", len(out.Mismatched))
	}
	if len(out.MissingInRight) != 0 || len(out.MissingInLeft) != 0 {
		t.Error("expected no missing entries for this key filter")
	}
}

func TestApply_KeyFilter_NoMatch(t *testing.T) {
	out := filter.Apply(baseResult(), filter.Options{Keys: []string{"NONEXISTENT"}})
	if len(out.MissingInRight) != 0 || len(out.MissingInLeft) != 0 || len(out.Mismatched) != 0 {
		t.Error("expected empty result for non-matching key filter")
	}
}

func TestApply_KeyFilter_MissingKey(t *testing.T) {
	out := filter.Apply(baseResult(), filter.Options{Keys: []string{"ONLY_LEFT"}})
	if len(out.MissingInRight) != 1 || out.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("expected ONLY_LEFT in MissingInRight, got %v", out.MissingInRight)
	}
}
