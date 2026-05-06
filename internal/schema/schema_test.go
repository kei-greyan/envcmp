package schema_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcmp/internal/schema"
)

func writeTempSchema(t *testing.T, s schema.Schema) string {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "schema.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadFile_ValidSchema(t *testing.T) {
	s := schema.Schema{
		"APP_ENV": {Required: true, Type: schema.TypeString},
	}
	path := writeTempSchema(t, s)
	loaded, err := schema.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := loaded["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in loaded schema")
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := schema.LoadFile("/nonexistent/schema.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFile_CorruptJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(path, []byte("not json{"), 0644)
	_, err := schema.LoadFile(path)
	if err == nil {
		t.Error("expected error for corrupt JSON")
	}
}

func TestValidate_NoViolations(t *testing.T) {
	env := map[string]string{"PORT": "8080", "DEBUG": "true"}
	s := schema.Schema{
		"PORT":  {Required: true, Type: schema.TypeInt},
		"DEBUG": {Required: true, Type: schema.TypeBool},
	}
	violations := schema.Validate(env, s)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	env := map[string]string{}
	s := schema.Schema{"SECRET": {Required: true, Type: schema.TypeString}}
	violations := schema.Validate(env, s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %q", violations[0].Key)
	}
}

func TestValidate_TypeMismatch_Int(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	s := schema.Schema{"PORT": {Type: schema.TypeInt}}
	violations := schema.Validate(env, s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_TypeMismatch_URL(t *testing.T) {
	env := map[string]string{"API_URL": "not-a-url"}
	s := schema.Schema{"API_URL": {Type: schema.TypeURL}}
	violations := schema.Validate(env, s)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	env := map[string]string{"VERSION": "abc"}
	s := schema.Schema{"VERSION": {Pattern: `^v\d+\.\d+`}}
	violations := schema.Validate(env, s)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_PatternMatches(t *testing.T) {
	env := map[string]string{"VERSION": "v1.2.3"}
	s := schema.Schema{"VERSION": {Pattern: `^v\d+\.\d+`}}
	violations := schema.Validate(env, s)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}
