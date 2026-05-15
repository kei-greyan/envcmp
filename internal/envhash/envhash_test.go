package envhash_test

import (
	"testing"

	"github.com/user/envcmp/internal/envhash"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"LOG_LEVEL": "info",
	}
}

func TestCompute_DeterministicHash(t *testing.T) {
	a := envhash.Compute(baseEnv())
	b := envhash.Compute(baseEnv())
	if a.Hash != b.Hash {
		t.Errorf("expected identical hashes, got %s vs %s", a.Hash, b.Hash)
	}
}

func TestCompute_PerKeyEntries(t *testing.T) {
	r := envhash.Compute(baseEnv())
	for _, k := range []string{"APP_ENV", "DB_HOST", "DB_PORT", "LOG_LEVEL"} {
		if _, ok := r.Entries[k]; !ok {
			t.Errorf("expected entry for key %q", k)
		}
	}
}

func TestCompute_EmptyEnv(t *testing.T) {
	r := envhash.Compute(map[string]string{})
	if r.Hash == "" {
		t.Error("expected non-empty hash for empty env")
	}
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestEqual_SameEnv_ReturnsTrue(t *testing.T) {
	a := envhash.Compute(baseEnv())
	b := envhash.Compute(baseEnv())
	if !envhash.Equal(a, b) {
		t.Error("expected Equal to return true for identical envs")
	}
}

func TestEqual_DifferentEnv_ReturnsFalse(t *testing.T) {
	env2 := baseEnv()
	env2["DB_HOST"] = "remotehost"
	a := envhash.Compute(baseEnv())
	b := envhash.Compute(env2)
	if envhash.Equal(a, b) {
		t.Error("expected Equal to return false for different envs")
	}
}

func TestDiff_ChangedKey_ReturnedInDiff(t *testing.T) {
	env2 := baseEnv()
	env2["DB_PORT"] = "5433"
	a := envhash.Compute(baseEnv())
	b := envhash.Compute(env2)
	diff := envhash.Diff(a, b)
	if len(diff) != 1 || diff[0] != "DB_PORT" {
		t.Errorf("expected [DB_PORT], got %v", diff)
	}
}

func TestDiff_AddedKey_IncludedInDiff(t *testing.T) {
	env2 := baseEnv()
	env2["NEW_KEY"] = "value"
	a := envhash.Compute(baseEnv())
	b := envhash.Compute(env2)
	diff := envhash.Diff(a, b)
	found := false
	for _, k := range diff {
		if k == "NEW_KEY" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected NEW_KEY in diff, got %v", diff)
	}
}

func TestDiff_NoDifference_ReturnsEmpty(t *testing.T) {
	a := envhash.Compute(baseEnv())
	b := envhash.Compute(baseEnv())
	diff := envhash.Diff(a, b)
	if len(diff) != 0 {
		t.Errorf("expected empty diff, got %v", diff)
	}
}

func TestDiff_IsSorted(t *testing.T) {
	env2 := baseEnv()
	env2["APP_ENV"] = "staging"
	env2["DB_HOST"] = "other"
	a := envhash.Compute(baseEnv())
	b := envhash.Compute(env2)
	diff := envhash.Diff(a, b)
	for i := 1; i < len(diff); i++ {
		if diff[i] < diff[i-1] {
			t.Errorf("diff not sorted at index %d: %v", i, diff)
		}
	}
}
