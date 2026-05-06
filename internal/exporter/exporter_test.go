package exporter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/exporter"
)

func baseResult() comparator.Result {
	return comparator.Result{
		MissingInRight: []string{"MISSING_KEY"},
		MissingInLeft:  []string{},
		Mismatched: []comparator.Diff{
			{Key: "APP_ENV", LeftVal: "dev", RightVal: "prod"},
		},
	}
}

func TestExport_ToStdout_NoError(t *testing.T) {
	opts := exporter.Options{Format: "text"}
	if err := exporter.Export(baseResult(), opts); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestExport_ToFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "result.txt")

	opts := exporter.Options{Format: "text", OutputFile: out, Overwrite: false}
	if err := exporter.Export(baseResult(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "MISSING_KEY") {
		t.Errorf("expected output to contain MISSING_KEY, got: %s", string(data))
	}
}

func TestExport_ToFile_JSON_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "result.json")

	opts := exporter.Options{Format: "json", OutputFile: out, Overwrite: false}
	if err := exporter.Export(baseResult(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "missing_in_right") {
		t.Errorf("expected JSON output, got: %s", string(data))
	}
}

func TestExport_NoOverwrite_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "result.txt")

	// Create the file first
	if err := os.WriteFile(out, []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := exporter.Options{Format: "text", OutputFile: out, Overwrite: false}
	err := exporter.Export(baseResult(), opts)
	if err == nil {
		t.Fatal("expected error when file exists and overwrite is false")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExport_Overwrite_ReplacesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "result.txt")

	if err := os.WriteFile(out, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := exporter.Options{Format: "text", OutputFile: out, Overwrite: true}
	if err := exporter.Export(baseResult(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "old content") {
		t.Errorf("expected file to be overwritten, still has old content")
	}
}
