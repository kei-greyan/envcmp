package envsplit_test

import (
	"testing"

	"github.com/user/envcmp/internal/envsplit"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"APP_NAME":    "envcmp",
		"APP_VERSION": "1.0.0",
		"AWS_KEY":     "abc",
		"UNRELATED":   "yes",
	}
}

func TestSplit_GroupsByPrefix(t *testing.T) {
	r := envsplit.Split(baseEnv(), []string{"DB_", "APP_"})

	dbGroup := findGroup(r.Groups, "DB_")
	if dbGroup == nil {
		t.Fatal("expected DB_ group")
	}
	if len(dbGroup.Keys) != 2 {
		t.Errorf("expected 2 DB_ keys, got %d", len(dbGroup.Keys))
	}

	appGroup := findGroup(r.Groups, "APP_")
	if appGroup == nil {
		t.Fatal("expected APP_ group")
	}
	if len(appGroup.Keys) != 2 {
		t.Errorf("expected 2 APP_ keys, got %d", len(appGroup.Keys))
	}
}

func TestSplit_UngroupedContainsUnmatched(t *testing.T) {
	r := envsplit.Split(baseEnv(), []string{"DB_", "APP_"})

	if _, ok := r.Ungrouped["AWS_KEY"]; !ok {
		t.Error("expected AWS_KEY in ungrouped")
	}
	if _, ok := r.Ungrouped["UNRELATED"]; !ok {
		t.Error("expected UNRELATED in ungrouped")
	}
	if len(r.Ungrouped) != 2 {
		t.Errorf("expected 2 ungrouped keys, got %d", len(r.Ungrouped))
	}
}

func TestSplit_LongestPrefixWins(t *testing.T) {
	env := map[string]string{
		"DB_REPLICA_HOST": "replica",
		"DB_HOST":         "primary",
	}
	r := envsplit.Split(env, []string{"DB_", "DB_REPLICA_"})

	replica := findGroup(r.Groups, "DB_REPLICA_")
	if replica == nil {
		t.Fatal("expected DB_REPLICA_ group")
	}
	if _, ok := replica.Keys["DB_REPLICA_HOST"]; !ok {
		t.Error("DB_REPLICA_HOST should be in DB_REPLICA_ group")
	}

	db := findGroup(r.Groups, "DB_")
	if db == nil {
		t.Fatal("expected DB_ group")
	}
	if _, ok := db.Keys["DB_HOST"]; !ok {
		t.Error("DB_HOST should be in DB_ group")
	}
	if _, ok := db.Keys["DB_REPLICA_HOST"]; ok {
		t.Error("DB_REPLICA_HOST should NOT be in DB_ group")
	}
}

func TestSplit_EmptyPrefixes_AllUngrouped(t *testing.T) {
	r := envsplit.Split(baseEnv(), []string{})
	if len(r.Ungrouped) != len(baseEnv()) {
		t.Errorf("expected all keys ungrouped, got %d", len(r.Ungrouped))
	}
	if len(r.Groups) != 0 {
		t.Errorf("expected no groups, got %d", len(r.Groups))
	}
}

func TestFormat_IncludesGroupHeaders(t *testing.T) {
	r := envsplit.Split(baseEnv(), []string{"DB_"})
	out := envsplit.Format(r)
	if out == "" {
		t.Fatal("expected non-empty format output")
	}
	if !contains(out, "[DB_]") {
		t.Error("expected [DB_] header in output")
	}
	if !contains(out, "[ungrouped]") {
		t.Error("expected [ungrouped] header in output")
	}
}

func findGroup(groups []envsplit.Group, prefix string) *envsplit.Group {
	for i := range groups {
		if groups[i].Prefix == prefix {
			return &groups[i]
		}
	}
	return nil
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
