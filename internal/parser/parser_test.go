package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", env["APP_ENV"])
	}
	if env["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %q", env["DEBUG"])
	}
}

func TestParseFile_CommentsAndEmptyLines(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
	if env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", env["KEY"])
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret"` + "\n" + `TOKEN='abc123'` + "\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "my secret" {
		t.Errorf("expected 'my secret', got %q", env["SECRET"])
	}
	if env["TOKEN"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", env["TOKEN"])
	}
}

func TestParseFile_KeyWithoutValue(t *testing.T) {
	path := writeTempEnv(t, "EMPTY_KEY=\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := env["EMPTY_KEY"]; !ok || v != "" {
		t.Errorf("expected EMPTY_KEY to exist with empty value, got %q", v)
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
