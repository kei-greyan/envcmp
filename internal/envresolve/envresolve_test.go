package envresolve_test

import (
	"testing"

	"github.com/user/envcmp/internal/envresolve"
)

func baseEnv() map[string]string {
	return map[string]string{
		"HOST":     "localhost",
		"PORT":     "5432",
		"DB_URL":   "postgres://${HOST}:${PORT}/mydb",
		"APP_NAME": "myapp",
	}
}

func TestResolve_ExpandsBracketSyntax(t *testing.T) {
	result, err := envresolve.Resolve(baseEnv(), envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://localhost:5432/mydb"
	if got := result.Resolved["DB_URL"]; got != want {
		t.Errorf("DB_URL: got %q, want %q", got, want)
	}
}

func TestResolve_NoBraces_DollarSyntax(t *testing.T) {
	env := map[string]string{
		"USER":    "alice",
		"GREETING": "hello $USER",
	}
	result, err := envresolve.Resolve(env, envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "hello alice"
	if got := result.Resolved["GREETING"]; got != want {
		t.Errorf("GREETING: got %q, want %q", got, want)
	}
}

func TestResolve_MissingRef_ProducesWarning(t *testing.T) {
	env := map[string]string{
		"VAL": "${UNDEFINED}",
	}
	opts := envresolve.DefaultOptions()
	opts.FailOnMissing = false
	result, err := envresolve.Resolve(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Warnings) == 0 {
		t.Error("expected at least one warning for undefined reference")
	}
	if got := result.Resolved["VAL"]; got != "" {
		t.Errorf("expected empty string for missing ref, got %q", got)
	}
}

func TestResolve_FailOnMissing_ReturnsError(t *testing.T) {
	env := map[string]string{
		"VAL": "${UNDEFINED}",
	}
	opts := envresolve.DefaultOptions()
	opts.FailOnMissing = true
	_, err := envresolve.Resolve(env, opts)
	if err == nil {
		t.Error("expected error for missing reference with FailOnMissing=true")
	}
}

func TestResolve_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	origDB := env["DB_URL"]
	_, err := envresolve.Resolve(env, envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB_URL"] != origDB {
		t.Error("Resolve mutated the original env map")
	}
}

func TestResolve_NoReferences_PassesThrough(t *testing.T) {
	env := map[string]string{
		"PLAIN": "just a value",
	}
	result, err := envresolve.Resolve(env, envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result.Resolved["PLAIN"]; got != "just a value" {
		t.Errorf("PLAIN: got %q, want %q", got, "just a value")
	}
}

func TestResolve_ChainedReferences(t *testing.T) {
	env := map[string]string{
		"A": "foo",
		"B": "${A}_bar",
		"C": "${B}_baz",
	}
	result, err := envresolve.Resolve(env, envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result.Resolved["C"]; got != "foo_bar_baz" {
		t.Errorf("C: got %q, want %q", got, "foo_bar_baz")
	}
}
