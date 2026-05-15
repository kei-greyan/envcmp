package envpin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"HOST":    "localhost",
		"PORT":    "8080",
		"DEBUG":   "true",
		"RATIO":   "1.5",
		"COMMENT": "",
	}
}

func TestCreate_InfersKinds(t *testing.T) {
	p := Create(baseEnv())
	cases := map[string]string{
		"HOST":  "string",
		"PORT":  "int",
		"DEBUG": "bool",
		"RATIO": "float",
	}
	for k, want := range cases {
		if got := p.Keys[k].Kind; got != want {
			t.Errorf("key %s: want kind %q, got %q", k, want, got)
		}
	}
}

func TestCreate_TracksNonEmpty(t *testing.T) {
	p := Create(baseEnv())
	if p.Keys["HOST"].NonEmpty != true {
		t.Error("HOST should be non-empty")
	}
	if p.Keys["COMMENT"].NonEmpty != false {
		t.Error("COMMENT should be empty")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	p := Create(baseEnv())
	path := filepath.Join(t.TempDir(), "pin.json")
	if err := SaveFile(path, p); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if len(loaded.Keys) != len(p.Keys) {
		t.Errorf("key count mismatch: want %d, got %d", len(p.Keys), len(loaded.Keys))
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := LoadFile("/nonexistent/pin.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pin.json")
	os.WriteFile(path, []byte("not json{"), 0644)
	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}

func TestCheck_NoViolations(t *testing.T) {
	p := Create(baseEnv())
	v := Check(p, baseEnv())
	if len(v) != 0 {
		t.Errorf("expected no violations, got %d", len(v))
	}
}

func TestCheck_MissingKey_ReturnsViolation(t *testing.T) {
	p := Create(baseEnv())
	env := baseEnv()
	delete(env, "PORT")
	v := Check(p, env)
	if len(v) != 1 || v[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %v", v)
	}
}

func TestCheck_TypeDrift_ReturnsViolation(t *testing.T) {
	p := Create(baseEnv())
	env := baseEnv()
	env["PORT"] = "not-a-number"
	v := Check(p, env)
	if len(v) != 1 || v[0].Key != "PORT" {
		t.Errorf("expected type drift violation for PORT, got %v", v)
	}
}

func TestCheck_EmptyValueViolation(t *testing.T) {
	p := Create(baseEnv())
	env := baseEnv()
	env["HOST"] = ""
	v := Check(p, env)
	if len(v) != 1 || v[0].Key != "HOST" {
		t.Errorf("expected non-empty violation for HOST, got %v", v)
	}
}

func TestSaveFile_InvalidPath_ReturnsError(t *testing.T) {
	p := Create(baseEnv())
	err := SaveFile("/no/such/dir/pin.json", p)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestCreate_EmptyEnv(t *testing.T) {
	p := Create(map[string]string{})
	data, _ := json.Marshal(p)
	if string(data) == "" {
		t.Error("expected non-empty JSON")
	}
	if len(p.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(p.Keys))
	}
}
