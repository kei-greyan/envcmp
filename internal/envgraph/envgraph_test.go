package envgraph_test

import (
	"testing"

	"github.com/user/envcmp/internal/envgraph"
)

func baseEnv() map[string]string {
	return map[string]string{
		"HOST":     "localhost",
		"PORT":     "5432",
		"DSN":      "postgres://$HOST:$PORT/db",
		"BASE_URL": "http://${HOST}:${PORT}",
		"APP_DSN":  "$DSN",
	}
}

func TestBuild_ExtractsDeps(t *testing.T) {
	g := envgraph.Build(baseEnv())

	if len(g.Nodes) != 5 {
		t.Fatalf("expected 5 nodes, got %d", len(g.Nodes))
	}
	dsn := g.Nodes["DSN"]
	if len(dsn.Deps) != 2 {
		t.Errorf("DSN: expected 2 deps, got %d", len(dsn.Deps))
	}
}

func TestBuild_NoDeps_EmptySlice(t *testing.T) {
	g := envgraph.Build(map[string]string{"KEY": "plain_value"})
	node := g.Nodes["KEY"]
	if len(node.Deps) != 0 {
		t.Errorf("expected no deps, got %v", node.Deps)
	}
}

func TestTopoSort_ReturnsAllKeys(t *testing.T) {
	g := envgraph.Build(baseEnv())
	order, err := g.TopoSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 5 {
		t.Errorf("expected 5 keys, got %d", len(order))
	}
}

func TestTopoSort_DepsBeforeDependents(t *testing.T) {
	g := envgraph.Build(baseEnv())
	order, err := g.TopoSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := make(map[string]int, len(order))
	for i, k := range order {
		pos[k] = i
	}
	if pos["HOST"] >= pos["DSN"] {
		t.Errorf("HOST should appear before DSN in topo order")
	}
	if pos["DSN"] >= pos["APP_DSN"] {
		t.Errorf("DSN should appear before APP_DSN in topo order")
	}
}

func TestTopoSort_CycleDetected(t *testing.T) {
	env := map[string]string{
		"A": "$B",
		"B": "$A",
	}
	g := envgraph.Build(env)
	_, err := g.TopoSort()
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestRoots_ReturnsKeysNobodyDependsOn(t *testing.T) {
	g := envgraph.Build(baseEnv())
	roots := g.Roots()
	// BASE_URL and APP_DSN are not referenced by anyone
	rootSet := make(map[string]bool)
	for _, r := range roots {
		rootSet[r] = true
	}
	if !rootSet["BASE_URL"] {
		t.Errorf("expected BASE_URL to be a root")
	}
	if !rootSet["APP_DSN"] {
		t.Errorf("expected APP_DSN to be a root")
	}
	if rootSet["HOST"] {
		t.Errorf("HOST is depended upon, should not be a root")
	}
}

func TestBuild_EmptyEnv_EmptyGraph(t *testing.T) {
	g := envgraph.Build(map[string]string{})
	if len(g.Nodes) != 0 {
		t.Errorf("expected empty graph")
	}
}
