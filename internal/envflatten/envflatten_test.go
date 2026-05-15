package envflatten_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/envflatten"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB__HOST":    "localhost",
		"DB__PORT":    "5432",
		"APP__DEBUG":  "true",
		"APP__NAME":   "envcmp",
		"STANDALONE":  "yes",
	}
}

func TestFlatten_GroupsByPrefix(t *testing.T) {
	groups := envflatten.Flatten(baseEnv(), envflatten.DefaultOptions())

	if _, ok := groups["DB"]; !ok {
		t.Fatal("expected group DB")
	}
	if _, ok := groups["APP"]; !ok {
		t.Fatal("expected group APP")
	}
	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 2 {
		t.Errorf("expected 2 APP entries, got %d", len(groups["APP"]))
	}
}

func TestFlatten_NoPrefixKey_PlacedUnderEmpty(t *testing.T) {
	groups := envflatten.Flatten(baseEnv(), envflatten.DefaultOptions())

	unprefixed, ok := groups[""]
	if !ok {
		t.Fatal("expected empty-prefix group for STANDALONE")
	}
	if len(unprefixed) != 1 || unprefixed[0].Key != "STANDALONE" {
		t.Errorf("unexpected unprefixed entries: %+v", unprefixed)
	}
}

func TestFlatten_EntriesAreSortedByKey(t *testing.T) {
	groups := envflatten.Flatten(baseEnv(), envflatten.DefaultOptions())

	db := groups["DB"]
	if db[0].Key != "DB__HOST" || db[1].Key != "DB__PORT" {
		t.Errorf("expected sorted DB keys, got %v, %v", db[0].Key, db[1].Key)
	}
}

func TestFlatten_Uppercase_ConvertsKeys(t *testing.T) {
	env := map[string]string{"db__host": "localhost", "db__port": "5432"}
	opts := envflatten.Options{Delimiter: "__", Uppercase: true}
	groups := envflatten.Flatten(env, opts)

	if _, ok := groups["DB"]; !ok {
		t.Fatal("expected uppercase group DB")
	}
}

func TestFlatten_CustomDelimiter(t *testing.T) {
	env := map[string]string{"DB.HOST": "localhost", "DB.PORT": "5432"}
	opts := envflatten.Options{Delimiter: "."}
	groups := envflatten.Flatten(env, opts)

	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 entries with custom delimiter, got %d", len(groups["DB"]))
	}
}

func TestFlatten_EmptyEnv_ReturnsEmptyGroups(t *testing.T) {
	groups := envflatten.Flatten(map[string]string{}, envflatten.DefaultOptions())
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestFormat_ContainsPrefixHeaders(t *testing.T) {
	groups := envflatten.Flatten(baseEnv(), envflatten.DefaultOptions())
	out := envflatten.Format(groups)

	for _, hdr := range []string{"[DB]", "[APP]", "[no prefix]"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormat_EmptyGroups_ReturnsEmptyMarker(t *testing.T) {
	out := envflatten.Format(map[string][]envflatten.FlatEntry{})
	if !strings.Contains(out, "(empty)") {
		t.Errorf("expected (empty) marker, got: %s", out)
	}
}

func TestSummary_ReturnsCorrectCounts(t *testing.T) {
	groups := envflatten.Flatten(baseEnv(), envflatten.DefaultOptions())
	s := envflatten.Summary(groups)

	if !strings.Contains(s, "5 key(s)") {
		t.Errorf("expected 5 keys in summary, got: %s", s)
	}
	if !strings.Contains(s, "3 group(s)") {
		t.Errorf("expected 3 groups in summary, got: %s", s)
	}
}
