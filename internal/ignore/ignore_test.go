package ignore_test

import (
	"os"
	"testing"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/ignore"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.ignore")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func baseResult() comparator.Result {
	return comparator.Result{
		MissingInRight: []string{"SECRET", "PORT"},
		MissingInLeft:  []string{"DB_HOST"},
		Mismatched: []comparator.Mismatch{
			{Key: "APP_ENV", LeftVal: "dev", RightVal: "prod"},
		},
	}
}

func TestLoadFile_BasicKeys(t *testing.T) {
	path := writeTempIgnore(t, "# ignore these\nSECRET\nPORT\n")
	keys, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := keys["SECRET"]; !ok {
		t.Error("expected SECRET in ignore set")
	}
	if _, ok := keys["PORT"]; !ok {
		t.Error("expected PORT in ignore set")
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestLoadFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempIgnore(t, "\n# comment\n  \nKEY_A\n")
	keys, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := ignore.LoadFile("/nonexistent/path/.ignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestApply_RemovesIgnoredKeys(t *testing.T) {
	ignoreKeys := map[string]struct{}{"SECRET": {}, "DB_HOST": {}, "APP_ENV": {}}
	result := ignore.Apply(baseResult(), ignoreKeys)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "PORT" {
		t.Errorf("expected only PORT in MissingInRight, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected empty MissingInLeft, got %v", result.MissingInLeft)
	}
	if len(result.Mismatched) != 0 {
		t.Errorf("expected empty Mismatched, got %v", result.Mismatched)
	}
}

func TestApply_EmptyIgnoreSet_ReturnsAll(t *testing.T) {
	result := ignore.Apply(baseResult(), map[string]struct{}{})
	if len(result.MissingInRight) != 2 {
		t.Errorf("expected 2 MissingInRight, got %d", len(result.MissingInRight))
	}
}
