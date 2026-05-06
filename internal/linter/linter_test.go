package linter_test

import (
	"strings"
	"testing"

	"github.com/user/envcmp/internal/linter"
)

func lines(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n")
}

func TestLint_CleanFile(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	ls := lines("HOST=localhost\nPORT=8080")
	r := linter.Lint("test.env", ls, env)
	if !r.IsClean() {
		t.Fatalf("expected no issues, got %+v", r.Issues)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	env := map[string]string{"HOST": "b"}
	ls := lines("HOST=a\nHOST=b")
	r := linter.Lint("test.env", ls, env)
	if r.IsClean() {
		t.Fatal("expected duplicate key issue")
	}
	found := false
	for _, iss := range r.Issues {
		if iss.Key == "HOST" && strings.Contains(iss.Message, "duplicate") {
			found = true
		}
	}
	if !found {
		t.Errorf("duplicate issue not reported: %+v", r.Issues)
	}
}

func TestLint_InvalidKeyName(t *testing.T) {
	env := map[string]string{"1INVALID": "val"}
	ls := lines("1INVALID=val")
	r := linter.Lint("test.env", ls, env)
	if r.IsClean() {
		t.Fatal("expected invalid key issue")
	}
	if !strings.Contains(r.Issues[0].Message, "does not match") {
		t.Errorf("unexpected message: %s", r.Issues[0].Message)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{"TOKEN": ""}
	ls := lines("TOKEN=")
	r := linter.Lint("test.env", ls, env)
	if r.IsClean() {
		t.Fatal("expected empty value issue")
	}
	if !strings.Contains(r.Issues[0].Message, "empty") {
		t.Errorf("unexpected message: %s", r.Issues[0].Message)
	}
}

func TestLint_SkipsCommentsAndBlanks(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	ls := lines("# comment\n\nKEY=val")
	r := linter.Lint("test.env", ls, env)
	if !r.IsClean() {
		t.Fatalf("expected no issues, got %+v", r.Issues)
	}
}

func TestLint_InvalidAssignment(t *testing.T) {
	env := map[string]string{}
	ls := lines("NOTANASSIGNMENT")
	r := linter.Lint("test.env", ls, env)
	if r.IsClean() {
		t.Fatal("expected invalid assignment issue")
	}
	if !strings.Contains(r.Issues[0].Message, "not a valid") {
		t.Errorf("unexpected message: %s", r.Issues[0].Message)
	}
}
