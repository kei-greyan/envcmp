package envcast_test

import (
	"testing"

	"github.com/user/envcmp/internal/envcast"
)

func baseEnv() map[string]string {
	return map[string]string{
		"PORT":    "8080",
		"RATIO":   "3.14",
		"ENABLED": "true",
		"NAME":    "myapp",
	}
}

func TestCast_AllKnownTypes(t *testing.T) {
	types := map[string]string{
		"PORT":    "int",
		"RATIO":   "float",
		"ENABLED": "bool",
		"NAME":    "string",
	}
	results, err := envcast.Cast(baseEnv(), types, envcast.CastOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("key %q: unexpected cast error: %v", r.Key, r.Err)
		}
	}
}

func TestCast_InvalidInt_NonStrict(t *testing.T) {
	env := map[string]string{"PORT": "notanint"}
	types := map[string]string{"PORT": "int"}
	results, err := envcast.Cast(env, types, envcast.CastOptions{Strict: false})
	if err != nil {
		t.Fatalf("expected no error in non-strict mode, got %v", err)
	}
	if len(results) != 1 || results[0].Err == nil {
		t.Error("expected a failure result for invalid int")
	}
}

func TestCast_InvalidInt_Strict(t *testing.T) {
	env := map[string]string{"PORT": "notanint"}
	types := map[string]string{"PORT": "int"}
	_, err := envcast.Cast(env, types, envcast.CastOptions{Strict: true})
	if err == nil {
		t.Error("expected error in strict mode")
	}
}

func TestCast_UnknownType_ReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	types := map[string]string{"KEY": "uuid"}
	results, err := envcast.Cast(env, types, envcast.CastOptions{})
	if err != nil {
		t.Fatalf("unexpected hard error: %v", err)
	}
	if len(results) == 0 || results[0].Err == nil {
		t.Error("expected failure result for unknown type")
	}
}

func TestCast_NoTypesMap_DefaultsToString(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	results, err := envcast.Cast(env, map[string]string{}, envcast.CastOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Kind != "string" {
		t.Errorf("expected kind string, got %q", results[0].Kind)
	}
}

func TestFailures_ReturnsOnlyErrors(t *testing.T) {
	env := map[string]string{"PORT": "bad", "NAME": "ok"}
	types := map[string]string{"PORT": "int", "NAME": "string"}
	results, _ := envcast.Cast(env, types, envcast.CastOptions{})
	fails := envcast.Failures(results)
	if len(fails) != 1 {
		t.Errorf("expected 1 failure, got %d", len(fails))
	}
	if fails[0].Key != "PORT" {
		t.Errorf("expected failure on PORT, got %q", fails[0].Key)
	}
}
