package envsign_test

import (
	"testing"

	"github.com/user/envcmp/internal/envsign"
)

var baseEnv = map[string]string{
	"APP_ENV": "production",
	"DB_HOST": "localhost",
	"SECRET":  "s3cr3t",
}

func TestSign_Deterministic(t *testing.T) {
	sig1 := envsign.Sign(baseEnv, "my-secret")
	sig2 := envsign.Sign(baseEnv, "my-secret")
	if sig1 != sig2 {
		t.Fatalf("expected deterministic signature, got %q vs %q", sig1, sig2)
	}
}

func TestSign_DifferentSecrets_DifferentSigs(t *testing.T) {
	sig1 := envsign.Sign(baseEnv, "secret-a")
	sig2 := envsign.Sign(baseEnv, "secret-b")
	if sig1 == sig2 {
		t.Fatal("expected different signatures for different secrets")
	}
}

func TestSign_DifferentEnvs_DifferentSigs(t *testing.T) {
	other := map[string]string{"APP_ENV": "staging"}
	sig1 := envsign.Sign(baseEnv, "key")
	sig2 := envsign.Sign(other, "key")
	if sig1 == sig2 {
		t.Fatal("expected different signatures for different env maps")
	}
}

func TestVerify_ValidSignature_ReturnsTrue(t *testing.T) {
	sig := envsign.Sign(baseEnv, "my-secret")
	if !envsign.Verify(baseEnv, "my-secret", sig) {
		t.Fatal("expected Verify to return true for valid signature")
	}
}

func TestVerify_TamperedEnv_ReturnsFalse(t *testing.T) {
	sig := envsign.Sign(baseEnv, "my-secret")
	tampered := map[string]string{
		"APP_ENV": "development",
		"DB_HOST": "localhost",
		"SECRET":  "s3cr3t",
	}
	if envsign.Verify(tampered, "my-secret", sig) {
		t.Fatal("expected Verify to return false for tampered env")
	}
}

func TestVerify_WrongSecret_ReturnsFalse(t *testing.T) {
	sig := envsign.Sign(baseEnv, "correct-secret")
	if envsign.Verify(baseEnv, "wrong-secret", sig) {
		t.Fatal("expected Verify to return false for wrong secret")
	}
}

func TestFingerprint_Deterministic(t *testing.T) {
	fp1 := envsign.Fingerprint(baseEnv)
	fp2 := envsign.Fingerprint(baseEnv)
	if fp1 != fp2 {
		t.Fatalf("expected deterministic fingerprint, got %q vs %q", fp1, fp2)
	}
}

func TestFingerprint_Length(t *testing.T) {
	fp := envsign.Fingerprint(baseEnv)
	if len(fp) != 12 {
		t.Fatalf("expected fingerprint length 12, got %d", len(fp))
	}
}

func TestFingerprint_EmptyEnv(t *testing.T) {
	fp := envsign.Fingerprint(map[string]string{})
	if len(fp) != 12 {
		t.Fatalf("expected fingerprint length 12 for empty env, got %d", len(fp))
	}
}
