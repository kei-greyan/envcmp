package templater_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envcmp/internal/templater"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"APP_PORT":    "8080",
}

func TestGenerate_KeysAreSorted(t *testing.T) {
	out := templater.Generate(sampleEnv, templater.Options{})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APP_PORT=") {
		t.Errorf("expected first key APP_PORT, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "DB_HOST=") {
		t.Errorf("expected second key DB_HOST, got %s", lines[1])
	}
}

func TestGenerate_ValuesAreEmpty_ByDefault(t *testing.T) {
	out := templater.Generate(sampleEnv, templater.Options{})
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 || parts[1] != "" {
			t.Errorf("expected empty value, got line: %q", line)
		}
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	out := templater.Generate(sampleEnv, templater.Options{Placeholder: "CHANGEME"})
	if !strings.Contains(out, "=CHANGEME") {
		t.Errorf("expected placeholder CHANGEME in output:\n%s", out)
	}
}

func TestGenerate_AddComments(t *testing.T) {
	out := templater.Generate(sampleEnv, templater.Options{AddComments: true})
	if !strings.Contains(out, "# APP_PORT") {
		t.Errorf("expected comment for APP_PORT in output:\n%s", out)
	}
}

func TestGenerate_EmptyEnv(t *testing.T) {
	out := templater.Generate(map[string]string{}, templater.Options{})
	if out != "" {
		t.Errorf("expected empty output for empty env, got %q", out)
	}
}

func TestGenerate_SingleEntry(t *testing.T) {
	env := map[string]string{"ONLY_KEY": "some_value"}
	out := templater.Generate(env, templater.Options{})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != "ONLY_KEY=" {
		t.Errorf("expected 'ONLY_KEY=', got %q", lines[0])
	}
}

func TestWriteFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.template")

	err := templater.WriteFile(path, sampleEnv, templater.Options{})
	if err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}
	if !strings.Contains(string(data), "APP_PORT=") {
		t.Errorf("expected APP_PORT= in file, got:\n%s", string(data))
	}
}

func TestWriteFile_InvalidPath_ReturnsError(t *testing.T) {
	err := templater.WriteFile("/nonexistent/dir/.env.template", sampleEnv, templater.Options{})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestWriteFile_ContentMatchesGenerate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.template")
	opts := templater.Options{Placeholder: "FILL_ME", AddComments: true}

	err := templater.WriteFile(path, sampleEnv, opts)
	if err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}

	expected := templater.Generate(sampleEnv, opts)
	if strings.TrimSpace(string(data)) != strings.TrimSpace(expected) {
		t.Errorf("file content does not match Generate output\ngot:\n%s\nwant:\n%s", string(data), expected)
	}
}
