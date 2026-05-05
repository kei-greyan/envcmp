package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envcmp")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestCLI_NoDiff_ExitZero(t *testing.T) {
	bin := buildBinary(t)
	f1 := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	f2 := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	cmd := exec.Command(bin, f1, f2)
	if err := cmd.Run(); err != nil {
		t.Errorf("expected exit 0, got: %v", err)
	}
}

func TestCLI_Diff_StrictExitOne(t *testing.T) {
	bin := buildBinary(t)
	f1 := writeTempEnv(t, "FOO=bar\n")
	f2 := writeTempEnv(t, "FOO=different\n")

	cmd := exec.Command(bin, "-strict", f1, f2)
	err := cmd.Run()
	if err == nil {
		t.Error("expected non-zero exit code with -strict flag and differences")
	}
}

func TestCLI_MissingArgs_ExitTwo(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Errorf("expected exit code 2, got %d", exitErr.ExitCode())
		}
	} else {
		t.Error("expected an exit error")
	}
}

func TestCLI_QuietFlag_NoOutput(t *testing.T) {
	bin := buildBinary(t)
	f1 := writeTempEnv(t, "FOO=bar\n")
	f2 := writeTempEnv(t, "FOO=other\n")

	cmd := exec.Command(bin, "-quiet", f1, f2)
	out, _ := cmd.CombinedOutput()
	if len(out) != 0 {
		t.Errorf("expected no output with -quiet, got: %s", out)
	}
}
