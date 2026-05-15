package envgraph

import (
	"fmt"
	"sort"
	"strings"
)

// Format returns a human-readable representation of the dependency graph.
func Format(g *Graph) string {
	if len(g.Nodes) == 0 {
		return "(empty graph)\n"
	}
	var sb strings.Builder
	keys := sortedKeys(g.Nodes)
	for _, k := range keys {
		node := g.Nodes[k]
		if len(node.Deps) == 0 {
			fmt.Fprintf(&sb, "  %s (no deps)\n", k)
		} else {
			sorted := make([]string, len(node.Deps))
			copy(sorted, node.Deps)
			sort.Strings(sorted)
			fmt.Fprintf(&sb, "  %s -> [%s]\n", k, strings.Join(sorted, ", "))
		}
	}
	return sb.String()
}

// Summary returns a one-line summary of the graph.
func Summary(g *Graph) string {
	total := len(g.Nodes)
	withDeps := 0
	for _, node := range g.Nodes {
		if len(node.Deps) > 0 {
			withDeps++
		}
	}
	return fmt.Sprintf("%d keys, %d with references", total, withDeps)
}
