package envnorm_test

import (
	"testing"

	"github.com/user/envcmp/internal/envnorm"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST": "  localhost  ",
		"app_port": "8080",
		"Debug":    " true ",
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	opts := envnorm.DefaultOptions()
	result := envnorm.Normalize(baseEnv(), opts)

	if result["APP_HOST"] != "localhost" {
		t.Errorf("expected trimmed value, got %q", result["APP_HOST"])
	}
	if result["Debug"] != "true" {
		t.Errorf("expected trimmed value, got %q", result["Debug"])
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	opts := envnorm.Options{UppercaseKeys: true, TrimValues: false}
	result := envnorm.Normalize(baseEnv(), opts)

	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT to exist after uppercasing")
	}
	if _, ok := result["DEBUG"]; !ok {
		t.Error("expected DEBUG to exist after uppercasing")
	}
}

func TestNormalize_LowercaseKeys(t *testing.T) {
	opts := envnorm.Options{LowercaseKeys: true, TrimValues: false}
	result := envnorm.Normalize(baseEnv(), opts)

	if _, ok := result["app_host"]; !ok {
		t.Error("expected app_host to exist after lowercasing")
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	opts := envnorm.Options{TrimValues: true, UppercaseKeys: true}
	_ = envnorm.Normalize(env, opts)

	if _, ok := env["APP_HOST"]; !ok {
		t.Error("original map should not be mutated")
	}
	if env["APP_HOST"] != "  localhost  " {
		t.Error("original value should not be modified")
	}
}

func TestDiff_DetectsChangedKeys(t *testing.T) {
	env := map[string]string{
		"KEY": "  value  ",
		"CLEAN": "ok",
	}
	opts := envnorm.DefaultOptions()
	changed := envnorm.Diff(env, opts)

	if _, ok := changed["KEY"]; !ok {
		t.Error("expected KEY to appear in diff")
	}
	if _, ok := changed["CLEAN"]; ok {
		t.Error("CLEAN should not appear in diff")
	}
}

func TestDiff_EmptyEnv(t *testing.T) {
	opts := envnorm.DefaultOptions()
	changed := envnorm.Diff(map[string]string{}, opts)
	if len(changed) != 0 {
		t.Errorf("expected empty diff, got %v", changed)
	}
}
