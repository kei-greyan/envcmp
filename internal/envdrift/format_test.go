package envdrift_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envdrift"
)

func TestFormat_ContainsKeyNames(t *testing.T) {
	r := envdrift.Result{
		Entries: []envdrift.Entry{
			{Key: "MY_KEY", Kind: envdrift.DriftAdded, Current: "hello"},
		},
	}
	out := envdrift.Format(r)
	if !strings.Contains(out, "MY_KEY") {
		t.Errorf("expected key name in output, got: %q", out)
	}
}

func TestFormat_ChangedShowsBaselineAndCurrent(t *testing.T) {
	r := envdrift.Result{
		Entries: []envdrift.Entry{
			{Key: "PORT", Kind: envdrift.DriftChanged, Baseline: "8080", Current: "9090"},
		},
	}
	out := envdrift.Format(r)
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Errorf("expected both baseline and current values in output: %q", out)
	}
}

func TestFormat_IncludesSummaryLine(t *testing.T) {
	r := envdrift.Result{
		Entries: []envdrift.Entry{
			{Key: "X", Kind: envdrift.DriftRemoved, Baseline: "y"},
		},
	}
	out := envdrift.Format(r)
	if !strings.Contains(out, "drift detected") {
		t.Errorf("expected summary line in output: %q", out)
	}
}

func TestFormat_MultipleEntries_AllPresent(t *testing.T) {
	r := envdrift.Result{
		Entries: []envdrift.Entry{
			{Key: "ALPHA", Kind: envdrift.DriftAdded, Current: "1"},
			{Key: "BETA", Kind: envdrift.DriftRemoved, Baseline: "2"},
			{Key: "GAMMA", Kind: envdrift.DriftChanged, Baseline: "old", Current: "new"},
		},
	}
	out := envdrift.Format(r)
	for _, key := range []string{"ALPHA", "BETA", "GAMMA"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in output", key)
		}
	}
}
