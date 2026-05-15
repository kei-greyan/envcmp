package envgroup_test

import (
	"testing"

	"github.com/user/envcmp/internal/envgroup"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_ENV":     "production",
	"APP_VERSION": "1.2.3",
	"STANDALONE":  "yes",
}

func TestGroup_DefaultOptions_GroupsByFirstSegment(t *testing.T) {
	r, err := envgroup.Group(baseEnv, envgroup.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Groups["DB"]; !ok {
		t.Error("expected group DB")
	}
	if _, ok := r.Groups["APP"]; !ok {
		t.Error("expected group APP")
	}
	if _, ok := r.Groups["other"]; !ok {
		t.Error("expected group other for STANDALONE")
	}
}

func TestGroup_UngroupedLabel_IsCustomizable(t *testing.T) {
	opts := envgroup.DefaultOptions()
	opts.UngroupedLabel = "misc"
	r, _ := envgroup.Group(baseEnv, opts)
	if _, ok := r.Groups["misc"]; !ok {
		t.Error("expected custom ungrouped label 'misc'")
	}
}

func TestGroup_EmptyDelimiter_ReturnsError(t *testing.T) {
	opts := envgroup.DefaultOptions()
	opts.Delimiter = ""
	_, err := envgroup.Group(baseEnv, opts)
	if err == nil {
		t.Error("expected error for empty delimiter")
	}
}

func TestGroup_EmptyEnv_ReturnsEmptyResult(t *testing.T) {
	r, err := envgroup.Group(map[string]string{}, envgroup.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(r.Groups))
	}
}

func TestGroup_OrderIsSorted(t *testing.T) {
	r, _ := envgroup.Group(baseEnv, envgroup.DefaultOptions())
	for i := 1; i < len(r.Order); i++ {
		if r.Order[i-1] > r.Order[i] {
			t.Errorf("order not sorted: %v", r.Order)
		}
	}
}

func TestGroup_MaxDepth2_UsesDoubleSegment(t *testing.T) {
	env := map[string]string{
		"AWS_S3_BUCKET": "my-bucket",
		"AWS_S3_REGION": "us-east-1",
		"AWS_EC2_TYPE":  "t3.micro",
	}
	opts := envgroup.DefaultOptions()
	opts.MaxDepth = 2
	r, err := envgroup.Group(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Groups["AWS_S3"]; !ok {
		t.Error("expected group AWS_S3")
	}
	if _, ok := r.Groups["AWS_EC2"]; !ok {
		t.Error("expected group AWS_EC2")
	}
}
