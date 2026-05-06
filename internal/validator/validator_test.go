package validator_test

import (
	"testing"

	"github.com/user/envcmp/internal/validator"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":    "production",
		"PORT":       "8080",
		"SECRET_KEY": "abc123",
	}
}

func TestValidate_NoRules_NoViolations(t *testing.T) {
	v := validator.Validate(baseEnv(), nil)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestValidate_RequiredKeyPresent(t *testing.T) {
	rules := []validator.Rule{{Key: "APP_ENV", Required: true}}
	v := validator.Validate(baseEnv(), rules)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got: %v", v)
	}
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	rules := []validator.Rule{{Key: "DATABASE_URL", Required: true}}
	v := validator.Validate(baseEnv(), rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected key %q", v[0].Key)
	}
}

func TestValidate_RequiredKeyEmptyValue(t *testing.T) {
	env := baseEnv()
	env["SECRET_KEY"] = "   "
	rules := []validator.Rule{{Key: "SECRET_KEY", Required: true}}
	v := validator.Validate(env, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestValidate_PatternMatch_Valid(t *testing.T) {
	rules := []validator.Rule{{Key: "PORT", Pattern: `^\d+$`}}
	v := validator.Validate(baseEnv(), rules)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got: %v", v)
	}
}

func TestValidate_PatternMatch_Invalid(t *testing.T) {
	env := baseEnv()
	env["PORT"] = "not-a-number"
	rules := []validator.Rule{{Key: "PORT", Pattern: `^\d+$`}}
	v := validator.Validate(env, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestValidate_InvalidPattern_ReturnsViolation(t *testing.T) {
	rules := []validator.Rule{{Key: "APP_ENV", Pattern: `[invalid`}}
	v := validator.Validate(baseEnv(), rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation for bad pattern, got %d", len(v))
	}
}

func TestValidate_OptionalMissingKey_NoViolation(t *testing.T) {
	rules := []validator.Rule{{Key: "OPTIONAL_KEY", Pattern: `^\w+$`}}
	v := validator.Validate(baseEnv(), rules)
	if len(v) != 0 {
		t.Fatalf("expected no violations for optional missing key, got %d", len(v))
	}
}

func TestViolation_String(t *testing.T) {
	v := validator.Violation{Key: "FOO", Value: "bar", Message: "some error"}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string from Violation.String()")
	}
}
