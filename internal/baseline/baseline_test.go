package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcmp/internal/baseline"
	"github.com/user/envcmp/internal/comparator"
)

func baseResult() comparator.Result {
	return comparator.Result{
		MissingInRight: []string{"ALPHA"},
		MissingInLeft:  []string{"BETA"},
		Mismatched: map[string][2]string{
			"PORT": {"8080", "9090"},
		},
		Matched: []string{"HOST"},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	result := baseResult()
	if err := baseline.Save(path, "left.env", "right.env", result); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	rec, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: unexpected error: %v", err)
	}

	if rec.LeftFile != "left.env" {
		t.Errorf("LeftFile: got %q, want %q", rec.LeftFile, "left.env")
	}
	if rec.RightFile != "right.env" {
		t.Errorf("RightFile: got %q, want %q", rec.RightFile, "right.env")
	}
	if len(rec.Result.MissingInRight) != 1 || rec.Result.MissingInRight[0] != "ALPHA" {
		t.Errorf("MissingInRight: got %v", rec.Result.MissingInRight)
	}
	if len(rec.Result.MissingInLeft) != 1 || rec.Result.MissingInLeft[0] != "BETA" {
		t.Errorf("MissingInLeft: got %v", rec.Result.MissingInLeft)
	}
	if v, ok := rec.Result.Mismatched["PORT"]; !ok || v[0] != "8080" || v[1] != "9090" {
		t.Errorf("Mismatched[PORT]: got %v", rec.Result.Mismatched)
	}
	if rec.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestSave_InvalidPath_ReturnsError(t *testing.T) {
	result := baseResult()
	err := baseline.Save("/nonexistent/dir/baseline.json", "l", "r", result)
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for corrupt JSON, got nil")
	}
}
