package patcher_test

import (
	"testing"

	"github.com/your-org/envcmp/internal/patcher"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.prod.internal",
		"LOG_LEVEL": "info",
	}
}

func TestApply_Overwrite_UpdatesExistingKey(t *testing.T) {
	result, err := patcher.Apply(baseEnv(), map[string]string{"LOG_LEVEL": "debug"}, patcher.Overwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["LOG_LEVEL"] != "debug" {
		t.Errorf("expected debug, got %s", result.Env["LOG_LEVEL"])
	}
	if len(result.Updated) != 1 || result.Updated[0] != "LOG_LEVEL" {
		t.Errorf("expected Updated=[LOG_LEVEL], got %v", result.Updated)
	}
}

func TestApply_Overwrite_AddsNewKey(t *testing.T) {
	result, err := patcher.Apply(baseEnv(), map[string]string{"NEW_KEY": "value"}, patcher.Overwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["NEW_KEY"] != "value" {
		t.Errorf("expected value, got %s", result.Env["NEW_KEY"])
	}
	if len(result.Added) != 1 || result.Added[0] != "NEW_KEY" {
		t.Errorf("expected Added=[NEW_KEY], got %v", result.Added)
	}
}

func TestApply_SkipExisting_PreservesOriginal(t *testing.T) {
	result, err := patcher.Apply(baseEnv(), map[string]string{"APP_ENV": "staging"}, patcher.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", result.Env["APP_ENV"])
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "APP_ENV" {
		t.Errorf("expected Skipped=[APP_ENV], got %v", result.Skipped)
	}
}

func TestApply_ErrorOnConflict_ReturnsError(t *testing.T) {
	_, err := patcher.Apply(baseEnv(), map[string]string{"DB_HOST": "localhost"}, patcher.ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := baseEnv()
	_, err := patcher.Apply(base, map[string]string{"APP_ENV": "staging"}, patcher.Overwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if base["APP_ENV"] != "production" {
		t.Errorf("base was mutated: got %s", base["APP_ENV"])
	}
}

func TestApply_InvalidKey_ReturnsError(t *testing.T) {
	_, err := patcher.Apply(baseEnv(), map[string]string{"INVALID KEY": "x"}, patcher.Overwrite)
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

func TestApply_NilBase_ReturnsError(t *testing.T) {
	_, err := patcher.Apply(nil, map[string]string{}, patcher.Overwrite)
	if err == nil {
		t.Fatal("expected error for nil base, got nil")
	}
}
