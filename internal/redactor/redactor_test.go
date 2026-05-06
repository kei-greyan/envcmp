package redactor_test

import (
	"testing"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/redactor"
)

func baseResult() comparator.Result {
	return comparator.Result{
		MissingInRight: []string{"PLAIN_KEY"},
		MissingInLeft:  []string{"OTHER_KEY"},
		Mismatched: []comparator.Diff{
			{Key: "DB_PASSWORD", LeftVal: "hunter2", RightVal: "s3cr3t"},
			{Key: "APP_NAME", LeftVal: "foo", RightVal: "bar"},
			{Key: "API_KEY_PROD", LeftVal: "abc123", RightVal: "xyz789"},
		},
	}
}

func TestApply_MasksSensitiveValues(t *testing.T) {
	result := redactor.Apply(baseResult())

	for _, d := range result.Mismatched {
		if redactor.IsSensitive(d.Key) {
			if d.LeftVal != "***" || d.RightVal != "***" {
				t.Errorf("key %q: expected masked values, got left=%q right=%q", d.Key, d.LeftVal, d.RightVal)
			}
		}
	}
}

func TestApply_PreservesNonSensitiveValues(t *testing.T) {
	result := redactor.Apply(baseResult())

	for _, d := range result.Mismatched {
		if d.Key == "APP_NAME" {
			if d.LeftVal != "foo" || d.RightVal != "bar" {
				t.Errorf("APP_NAME values should not be masked, got left=%q right=%q", d.LeftVal, d.RightVal)
			}
		}
	}
}

func TestApply_DoesNotMutateMissingKeys(t *testing.T) {
	original := baseResult()
	result := redactor.Apply(original)

	if len(result.MissingInRight) != len(original.MissingInRight) {
		t.Fatalf("MissingInRight length changed")
	}
	if result.MissingInRight[0] != "PLAIN_KEY" {
		t.Errorf("expected PLAIN_KEY, got %q", result.MissingInRight[0])
	}
}

func TestApply_DoesNotModifyOriginal(t *testing.T) {
	original := baseResult()
	redactor.Apply(original)

	if original.Mismatched[0].LeftVal != "hunter2" {
		t.Error("Apply must not modify the original result")
	}
}

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"AUTH_TOKEN", true},
		{"api_key", true},
		{"PRIVATE_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"DEBUG", false},
	}
	for _, tc := range cases {
		got := redactor.IsSensitive(tc.key)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
