package watchrunner_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/envcmp/internal/watchrunner"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

func TestRun_NilWriter_ReturnsError(t *testing.T) {
	done := make(chan struct{})
	close(done)

	cfg := watchrunner.Config{
		LeftFile:  "left.env",
		RightFile: "right.env",
		Format:    "text",
		Out:       nil,
	}

	err := watchrunner.Run(cfg, done)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestRun_InitialComparisonPrinted(t *testing.T) {
	dir := t.TempDir()
	left := writeTempEnv(t, dir, "left.env", "KEY=hello\n")
	right := writeTempEnv(t, dir, "right.env", "KEY=hello\n")

	var buf bytes.Buffer
	done := make(chan struct{})

	cfg := watchrunner.Config{
		LeftFile:     left,
		RightFile:    right,
		Format:       "text",
		PollInterval: 50 * time.Millisecond,
		Out:          &buf,
	}

	go func() {
		time.Sleep(80 * time.Millisecond)
		close(done)
	}()

	_ = watchrunner.Run(cfg, done)

	if !strings.Contains(buf.String(), "comparison updated") {
		t.Errorf("expected header in output, got: %s", buf.String())
	}
}

func TestRun_RerunsOnFileChange(t *testing.T) {
	dir := t.TempDir()
	left := writeTempEnv(t, dir, "left.env", "KEY=1\n")
	right := writeTempEnv(t, dir, "right.env", "KEY=1\n")

	var buf bytes.Buffer
	done := make(chan struct{})

	cfg := watchrunner.Config{
		LeftFile:     left,
		RightFile:    right,
		Format:       "text",
		PollInterval: 40 * time.Millisecond,
		Out:          &buf,
	}

	go func() {
		_ = watchrunner.Run(cfg, done)
	}()

	time.Sleep(60 * time.Millisecond)
	_ = os.WriteFile(right, []byte("KEY=2\n"), 0o644)
	time.Sleep(150 * time.Millisecond)
	close(done)

	count := strings.Count(buf.String(), "comparison updated")
	if count < 2 {
		t.Errorf("expected at least 2 comparison prints, got %d\noutput:\n%s", count, buf.String())
	}
}

func TestRun_MissingFile_PrintsError(t *testing.T) {
	dir := t.TempDir()
	right := writeTempEnv(t, dir, "right.env", "KEY=1\n")

	var buf bytes.Buffer
	done := make(chan struct{})
	close(done)

	cfg := watchrunner.Config{
		LeftFile:  "/nonexistent/left.env",
		RightFile: right,
		Format:    "text",
		Out:       &buf,
	}

	// Should return an error because left file does not exist.
	err := watchrunner.Run(cfg, done)
	if err == nil {
		t.Error("expected error for missing left file")
	}
}
