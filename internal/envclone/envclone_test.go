package envclone_test

import (
	"testing"

	"github.com/user/envcmp/internal/envclone"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
	}
}

func TestClone_NoOptions_CopiesAll(t *testing.T) {
	res, err := envclone.Clone(baseEnv(), envclone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(res.Env))
	}
	if res.Env["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", res.Env["APP_HOST"])
	}
}

func TestClone_StripPrefix_RemovesPrefix(t *testing.T) {
	res, err := envclone.Clone(baseEnv(), envclone.Options{StripPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["HOST"]; !ok {
		t.Error("expected key HOST after stripping APP_ prefix")
	}
	if _, ok := res.Env["APP_HOST"]; ok {
		t.Error("original key APP_HOST should not be present")
	}
}

func TestClone_AddPrefix_InjectsPrefix(t *testing.T) {
	res, err := envclone.Clone(map[string]string{"HOST": "localhost"}, envclone.Options{AddPrefix: "SVC_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %q", res.Env["SVC_HOST"])
	}
}

func TestClone_OnlyKeys_FiltersKeys(t *testing.T) {
	opts := envclone.Options{
		OnlyKeys: map[string]struct{}{"APP_HOST": {}},
	}
	res, err := envclone.Clone(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 1 {
		t.Fatalf("expected 1 key, got %d", len(res.Env))
	}
	skipped := 0
	for _, e := range res.Entries {
		if e.Skipped {
			skipped++
		}
	}
	if skipped != 2 {
		t.Errorf("expected 2 skipped entries, got %d", skipped)
	}
}

func TestClone_KeyCollision_ReturnsError(t *testing.T) {
	src := map[string]string{
		"APP_HOST": "a",
		"HOST":     "b",
	}
	_, err := envclone.Clone(src, envclone.Options{StripPrefix: "APP_"})
	if err == nil {
		t.Fatal("expected collision error, got nil")
	}
}

func TestClone_DoesNotMutateOriginal(t *testing.T) {
	orig := baseEnv()
	_, _ = envclone.Clone(orig, envclone.Options{StripPrefix: "APP_", AddPrefix: "NEW_"})
	if _, ok := orig["APP_HOST"]; !ok {
		t.Error("original map was mutated")
	}
}
