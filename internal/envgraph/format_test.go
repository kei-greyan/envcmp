package envgraph_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envgraph"
)

func TestFormat_EmptyGraph_ReturnsPlaceholder(t *testing.T) {
	g := envgraph.Build(map[string]string{})
	out := envgraph.Format(g)
	if !strings.Contains(out, "empty") {
		t.Errorf("expected empty placeholder, got: %q", out)
	}
}

func TestFormat_ContainsKeyNames(t *testing.T) {
	g := envgraph.Build(map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://$HOST/db",
	})
	out := envgraph.Format(g)
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output")
	}
	if !strings.Contains(out, "DSN") {
		t.Errorf("expected DSN in output")
	}
}

func TestFormat_ShowsDepsInBrackets(t *testing.T) {
	g := envgraph.Build(map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://$HOST/db",
	})
	out := envgraph.Format(g)
	if !strings.Contains(out, "[HOST]") {
		t.Errorf("expected [HOST] in DSN line, got: %q", out)
	}
}

func TestFormat_NoDeps_ShowsNoDepsLabel(t *testing.T) {
	g := envgraph.Build(map[string]string{"PLAIN": "value"})
	out := envgraph.Format(g)
	if !strings.Contains(out, "no deps") {
		t.Errorf("expected 'no deps' label, got: %q", out)
	}
}

func TestSummary_CountsKeysAndRefs(t *testing.T) {
	g := envgraph.Build(map[string]string{
		"A": "plain",
		"B": "$A",
		"C": "${A}_${B}",
	})
	s := envgraph.Summary(g)
	if !strings.Contains(s, "3 keys") {
		t.Errorf("expected '3 keys' in summary, got: %q", s)
	}
	if !strings.Contains(s, "2 with references") {
		t.Errorf("expected '2 with references' in summary, got: %q", s)
	}
}
