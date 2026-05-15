package envgroup_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envgroup"
)

func TestFormat_ContainsGroupHeaders(t *testing.T) {
	r, _ := envgroup.Group(baseEnv, envgroup.DefaultOptions())
	out := envgroup.Format(r)
	if !strings.Contains(out, "[DB]") {
		t.Error("expected [DB] header in output")
	}
	if !strings.Contains(out, "[APP]") {
		t.Error("expected [APP] header in output")
	}
}

func TestFormat_ContainsKeyValues(t *testing.T) {
	r, _ := envgroup.Group(baseEnv, envgroup.DefaultOptions())
	out := envgroup.Format(r)
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Error("expected DB_HOST=localhost in output")
	}
}

func TestFormat_EmptyResult_ReturnsNoKeysMessage(t *testing.T) {
	r, _ := envgroup.Group(map[string]string{}, envgroup.DefaultOptions())
	out := envgroup.Format(r)
	if out != "(no keys)" {
		t.Errorf("unexpected output for empty result: %q", out)
	}
}

func TestSummary_CountsGroupsAndKeys(t *testing.T) {
	r, _ := envgroup.Group(baseEnv, envgroup.DefaultOptions())
	s := envgroup.Summary(r)
	if !strings.Contains(s, "group") {
		t.Error("expected 'group' in summary")
	}
	if !strings.Contains(s, "key") {
		t.Error("expected 'key' in summary")
	}
}

func TestSummary_EmptyResult(t *testing.T) {
	r, _ := envgroup.Group(map[string]string{}, envgroup.DefaultOptions())
	s := envgroup.Summary(r)
	if !strings.Contains(s, "0 group") {
		t.Errorf("expected '0 group' in summary, got: %s", s)
	}
}
