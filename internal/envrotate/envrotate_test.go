package envrotate

import (
	"strings"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "secret123",
		"API_KEY":     "abc",
		"APP_NAME":    "envcmp",
	}
}

func TestRotate_RandomStrategy_ReplacesValues(t *testing.T) {
	env := baseEnv()
	opts := DefaultOptions()
	opts.Keys = []string{"DB_PASSWORD", "API_KEY"}

	r, err := Rotate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["DB_PASSWORD"] == "secret123" {
		t.Error("expected DB_PASSWORD to be rotated")
	}
	if r.Env["API_KEY"] == "abc" {
		t.Error("expected API_KEY to be rotated")
	}
	if r.Env["APP_NAME"] != "envcmp" {
		t.Error("APP_NAME should not be modified")
	}
	if len(r.Rotated) != 2 {
		t.Errorf("expected 2 rotated entries, got %d", len(r.Rotated))
	}
}

func TestRotate_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	opts := DefaultOptions()
	opts.Keys = []string{"DB_PASSWORD"}

	_, err := Rotate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB_PASSWORD"] != "secret123" {
		t.Error("original map was mutated")
	}
}

func TestRotate_BlankStrategy_ClearsValues(t *testing.T) {
	env := baseEnv()
	opts := Options{Strategy: StrategyBlank, Keys: []string{"API_KEY"}, ByteLen: 16}

	r, err := Rotate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["API_KEY"] != "" {
		t.Errorf("expected empty value, got %q", r.Env["API_KEY"])
	}
}

func TestRotate_PlaceholderStrategy_SetsPlaceholder(t *testing.T) {
	env := baseEnv()
	opts := Options{Strategy: StrategyPlaceholder, Keys: []string{"DB_PASSWORD"}, ByteLen: 16}

	r, err := Rotate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["DB_PASSWORD"] != "ROTATED_DB_PASSWORD" {
		t.Errorf("unexpected placeholder: %q", r.Env["DB_PASSWORD"])
	}
}

func TestRotate_NoKeys_ReturnsError(t *testing.T) {
	_, err := Rotate(baseEnv(), DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty key list")
	}
}

func TestRotate_UnknownStrategy_ReturnsError(t *testing.T) {
	opts := Options{Strategy: "unknown", Keys: []string{"API_KEY"}, ByteLen: 16}
	_, err := Rotate(baseEnv(), opts)
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestRotate_RotatedLog_IsSortedByKey(t *testing.T) {
	env := baseEnv()
	opts := DefaultOptions()
	opts.Keys = []string{"DB_PASSWORD", "API_KEY"}

	r, err := Rotate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Rotated[0].Key != "API_KEY" || r.Rotated[1].Key != "DB_PASSWORD" {
		t.Errorf("expected sorted log, got %v, %v", r.Rotated[0].Key, r.Rotated[1].Key)
	}
}

func TestFormat_IncludesKeyNames(t *testing.T) {
	r := Result{
		Rotated: []Entry{
			{Key: "API_KEY", OldValue: "abc", NewValue: "xyz"},
		},
	}
	out := Format(r)
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in format output")
	}
}

func TestSummary_ReflectsCount(t *testing.T) {
	r := Result{
		Rotated: []Entry{
			{Key: "A"}, {Key: "B"},
		},
	}
	s := Summary(r)
	if !strings.Contains(s, "2 key(s)") {
		t.Errorf("unexpected summary: %q", s)
	}
}
