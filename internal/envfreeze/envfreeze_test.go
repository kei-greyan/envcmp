package envfreeze_test

import (
	"testing"

	"github.com/user/envcmp/internal/envfreeze"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"LOG_LEVEL": "info",
	}
}

func TestFreeze_SnapshotIsolated(t *testing.T) {
	env := baseEnv()
	s := envfreeze.Freeze(env)
	env["APP_ENV"] = "staging" // mutate original

	v, ok := s.Get("APP_ENV")
	if !ok || v != "production" {
		t.Errorf("expected frozen value %q, got %q", "production", v)
	}
}

func TestFreeze_Keys_Sorted(t *testing.T) {
	s := envfreeze.Freeze(baseEnv())
	keys := s.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestCheck_NoMutations_ReturnsEmpty(t *testing.T) {
	env := baseEnv()
	s := envfreeze.Freeze(env)
	violations := envfreeze.Check(s, baseEnv())
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_ModifiedKey_Detected(t *testing.T) {
	s := envfreeze.Freeze(baseEnv())
	current := baseEnv()
	current["APP_ENV"] = "staging"

	violations := envfreeze.Check(s, current)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Kind != "modified" || violations[0].Key != "APP_ENV" {
		t.Errorf("unexpected violation: %+v", violations[0])
	}
}

func TestCheck_AddedKey_Detected(t *testing.T) {
	s := envfreeze.Freeze(baseEnv())
	current := baseEnv()
	current["NEW_KEY"] = "value"

	violations := envfreeze.Check(s, current)
	if len(violations) != 1 || violations[0].Kind != "added" {
		t.Errorf("expected 1 added violation, got %+v", violations)
	}
}

func TestCheck_RemovedKey_Detected(t *testing.T) {
	s := envfreeze.Freeze(baseEnv())
	current := baseEnv()
	delete(current, "DB_HOST")

	violations := envfreeze.Check(s, current)
	if len(violations) != 1 || violations[0].Kind != "removed" {
		t.Errorf("expected 1 removed violation, got %+v", violations)
	}
}

func TestCheck_ViolationsSortedByKey(t *testing.T) {
	s := envfreeze.Freeze(baseEnv())
	current := map[string]string{"ZKEY": "z", "AKEY": "a"}

	violations := envfreeze.Check(s, current)
	for i := 1; i < len(violations); i++ {
		if violations[i].Key < violations[i-1].Key {
			t.Errorf("violations not sorted by key")
		}
	}
}

func TestSummary_Clean(t *testing.T) {
	msg := envfreeze.Summary(nil)
	if msg == "" {
		t.Error("expected non-empty summary")
	}
}

func TestFormat_NoViolations_ContainsNoMutations(t *testing.T) {
	out := envfreeze.Format(nil)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}

func TestFormat_WithViolations_ContainsKind(t *testing.T) {
	s := envfreeze.Freeze(baseEnv())
	current := baseEnv()
	current["APP_ENV"] = "dev"
	violations := envfreeze.Check(s, current)

	out := envfreeze.Format(violations)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
