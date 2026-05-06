// Package exporter handles writing formatted comparison results to an output
// destination, which can be standard output or a file on disk.
//
// It wraps the formatter package and adds file-system concerns such as
// destination path resolution, existence checks, and overwrite protection.
//
// Usage:
//
//	err := exporter.Export(result, exporter.Options{
//		Format:     "json",
//		OutputFile: "diff.json",
//		Overwrite:  false,
//	})
package exporter
