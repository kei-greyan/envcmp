package watcher_test

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/envcmp/internal/watcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestWatch_InvalidPaths_ReturnsError(t *testing.T) {
	done := make(chan struct{})
	close(done)

	err := watcher.Watch("", "right.env", watcher.DefaultOptions(), func() {}, done)
	if err == nil {
		t.Fatal("expected error for empty left path")
	}
}

func TestWatch_MissingFile_ReturnsError(t *testing.T) {
	done := make(chan struct{})
	close(done)

	err := watcher.Watch("/nonexistent/left.env", "/nonexistent/right.env", watcher.DefaultOptions(), func() {}, done)
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	left := writeTempEnv(t, "KEY=1\n")
	right := writeTempEnv(t, "KEY=1\n")

	var callCount atomic.Int32
	done := make(chan struct{})

	opts := watcher.Options{PollInterval: 50 * time.Millisecond}

	go func() {
		_ = watcher.Watch(left, right, opts, func() {
			callCount.Add(1)
		}, done)
	}()

	// Allow watcher to initialise.
	time.Sleep(80 * time.Millisecond)

	// Modify the left file.
	if err := os.WriteFile(left, []byte("KEY=2\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	time.Sleep(150 * time.Millisecond)
	close(done)

	if callCount.Load() == 0 {
		t.Error("expected onChange to be called at least once")
	}
}

func TestWatch_NoChange_CallbackNotInvoked(t *testing.T) {
	left := writeTempEnv(t, "KEY=1\n")
	right := writeTempEnv(t, "KEY=1\n")

	var callCount atomic.Int32
	done := make(chan struct{})

	opts := watcher.Options{PollInterval: 50 * time.Millisecond}

	go func() {
		_ = watcher.Watch(left, right, opts, func() {
			callCount.Add(1)
		}, done)
	}()

	time.Sleep(200 * time.Millisecond)
	close(done)

	if callCount.Load() != 0 {
		t.Errorf("expected no callback calls, got %d", callCount.Load())
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := watcher.DefaultOptions()
	if opts.PollInterval <= 0 {
		t.Error("expected positive poll interval")
	}
}

func TestWatch_StopsOnDone(t *testing.T) {
	left := writeTempEnv(t, "A=1\n")
	right := filepath.Join(filepath.Dir(left), "right.env")
	_ = os.WriteFile(right, []byte("A=1\n"), 0o644)

	done := make(chan struct{})
	errCh := make(chan error, 1)

	go func() {
		errCh <- watcher.Watch(left, right, watcher.Options{PollInterval: 30 * time.Millisecond}, func() {}, done)
	}()

	time.Sleep(60 * time.Millisecond)
	close(done)

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("watcher did not stop after done was closed")
	}
}
