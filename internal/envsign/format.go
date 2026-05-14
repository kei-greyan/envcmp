package envsign

import "fmt"

// VerifyResult holds the outcome of a signature verification.
type VerifyResult struct {
	Valid       bool
	Fingerprint string
	Signature   string
}

// Format returns a human-readable representation of a VerifyResult.
func Format(r VerifyResult) string {
	status := "VALID"
	if !r.Valid {
		status = "INVALID"
	}
	return fmt.Sprintf("signature: %s  fingerprint: %s  status: %s",
		r.Signature[:12]+"...", r.Fingerprint, status)
}

// Summary returns a one-line pass/fail summary.
func Summary(r VerifyResult) string {
	if r.Valid {
		return fmt.Sprintf("env signature OK (fingerprint %s)", r.Fingerprint)
	}
	return fmt.Sprintf("env signature MISMATCH (fingerprint %s)", r.Fingerprint)
}
