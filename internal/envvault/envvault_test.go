package envvault_test

import (
	"strings"
	"testing"

	"github.com/your-org/envcmp/internal/envvault"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"APP_NAME":    "envcmp",
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	env := baseEnv()
	enc, err := envvault.Encrypt(env, "passphrase")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	for k, v := range enc {
		if !strings.HasPrefix(v, "vault:") {
			t.Errorf("key %q: expected vault: prefix, got %q", k, v)
		}
	}
	dec, err := envvault.Decrypt(enc, "passphrase")
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	for k, want := range env {
		if got := dec[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestEncrypt_WrongPassphrase_FailsDecrypt(t *testing.T) {
	enc, err := envvault.Encrypt(baseEnv(), "correct")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	_, err = envvault.Decrypt(enc, "wrong")
	if err == nil {
		t.Fatal("expected error decrypting with wrong passphrase, got nil")
	}
}

func TestEncrypt_EmptyPassphrase_ReturnsError(t *testing.T) {
	_, err := envvault.Encrypt(baseEnv(), "")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_EmptyPassphrase_ReturnsError(t *testing.T) {
	_, err := envvault.Decrypt(baseEnv(), "")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestEncrypt_AlreadyEncrypted_SkipsValue(t *testing.T) {
	env := map[string]string{"KEY": "vault:alreadydone"}
	enc, err := envvault.Encrypt(env, "pass")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if enc["KEY"] != "vault:alreadydone" {
		t.Errorf("expected already-encrypted value to be unchanged, got %q", enc["KEY"])
	}
}

func TestDecrypt_PlaintextPassthrough(t *testing.T) {
	env := map[string]string{"APP_NAME": "envcmp"}
	dec, err := envvault.Decrypt(env, "pass")
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec["APP_NAME"] != "envcmp" {
		t.Errorf("expected plaintext passthrough, got %q", dec["APP_NAME"])
	}
}

func TestIsEncrypted_AllVault_ReturnsTrue(t *testing.T) {
	enc, _ := envvault.Encrypt(baseEnv(), "pass")
	if !envvault.IsEncrypted(enc) {
		t.Error("expected IsEncrypted to return true for fully encrypted env")
	}
}

func TestIsEncrypted_Mixed_ReturnsFalse(t *testing.T) {
	env := map[string]string{
		"A": "vault:abc",
		"B": "plaintext",
	}
	if envvault.IsEncrypted(env) {
		t.Error("expected IsEncrypted to return false for partially encrypted env")
	}
}

func TestEncrypt_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	orig := env["DB_PASSWORD"]
	_, err := envvault.Encrypt(env, "pass")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if env["DB_PASSWORD"] != orig {
		t.Error("Encrypt mutated the original map")
	}
}
