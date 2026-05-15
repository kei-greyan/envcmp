// Package envgraph analyses an env map and builds a directed dependency graph
// based on variable references found in values (e.g. $VAR or ${VAR}).
//
// It supports:
//   - Building a graph from any map[string]string
//   - Topological sorting for resolution order
//   - Cycle detection
//   - Identifying root keys (keys no other key depends on)
//   - Human-readable formatting
package envgraph
