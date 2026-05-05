package config

import (
	"testing"
)

func baseConfig() *Config {
	return &Config{
		LeftFile:  "a.env",
		RightFile: "b.env",
		Format:    FormatText,
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := baseConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_EmptyLeftFile(t *testing.T) {
	cfg := baseConfig()
	cfg.LeftFile = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty LeftFile")
	}
}

func TestValidate_EmptyRightFile(t *testing.T) {
	cfg := baseConfig()
	cfg.RightFile = "   "
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for blank RightFile")
	}
}

func TestValidate_DefaultFormatApplied(t *testing.T) {
	cfg := baseConfig()
	cfg.Format = ""
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != FormatText {
		t.Errorf("expected default format %q, got %q", FormatText, cfg.Format)
	}
}

func TestValidate_UnsupportedFormat(t *testing.T) {
	cfg := baseConfig()
	cfg.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestValidate_MutuallyExclusiveFlags(t *testing.T) {
	cfg := baseConfig()
	cfg.OnlyMissing = true
	cfg.OnlyMismatched = true
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when both only-missing and only-mismatched are set")
	}
}

func TestValidate_OnlyMissingAlone(t *testing.T) {
	cfg := baseConfig()
	cfg.OnlyMissing = true
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_JSONFormat(t *testing.T) {
	cfg := baseConfig()
	cfg.Format = FormatJSON
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
