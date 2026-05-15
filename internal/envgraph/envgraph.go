// Package envgraph builds a dependency graph from env variable references
// and detects cycles or resolution ordering.
package envgraph

import (
	"fmt"
	"sort"
	"strings"
)

// Node represents a single env variable and its dependencies.
type Node struct {
	Key  string
	Deps []string
}

// Graph holds the full dependency structure for an env map.
type Graph struct {
	Nodes map[string]*Node
}

// Build constructs a Graph from an env map by scanning values for $VAR or ${VAR} references.
func Build(env map[string]string) *Graph {
	g := &Graph{Nodes: make(map[string]*Node, len(env))}
	for k, v := range env {
		g.Nodes[k] = &Node{Key: k, Deps: extractRefs(v)}
	}
	return g
}

// TopoSort returns keys in dependency-resolved order.
// Returns an error if a cycle is detected.
func (g *Graph) TopoSort() ([]string, error) {
	visited := make(map[string]int) // 0=unvisited,1=visiting,2=done
	var order []string

	var visit func(key string) error
	visit = func(key string) error {
		switch visited[key] {
		case 2:
			return nil
		case 1:
			return fmt.Errorf("cycle detected at key: %s", key)
		}
		visited[key] = 1
		if node, ok := g.Nodes[key]; ok {
			for _, dep := range node.Deps {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		visited[key] = 2
		order = append(order, key)
		return nil
	}

	keys := sortedKeys(g.Nodes)
	for _, k := range keys {
		if err := visit(k); err != nil {
			return nil, err
		}
	}
	return order, nil
}

// Roots returns keys that no other key depends on.
func (g *Graph) Roots() []string {
	depended := make(map[string]bool)
	for _, node := range g.Nodes {
		for _, dep := range node.Deps {
			depended[dep] = true
		}
	}
	var roots []string
	for k := range g.Nodes {
		if !depended[k] {
			roots = append(roots, k)
		}
	}
	sort.Strings(roots)
	return roots
}

func extractRefs(value string) []string {
	var refs []string
	seen := make(map[string]bool)
	s := value
	for {
		start := strings.Index(s, "$")
		if start == -1 {
			break
		}
		s = s[start+1:]
		var key string
		if strings.HasPrefix(s, "{") {
			end := strings.Index(s, "}")
			if end == -1 {
				break
			}
			key = s[1:end]
			s = s[end+1:]
		} else {
			end := strings.IndexAny(s, " \t/:\n$")
			if end == -1 {
				key = s
				s = ""
			} else {
				key = s[:end]
				s = s[end:]
			}
		}
		if key != "" && !seen[key] {
			seen[key] = true
			refs = append(refs, key)
		}
	}
	return refs
}

func sortedKeys(m map[string]*Node) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
