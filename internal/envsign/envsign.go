// Package envsign provides HMAC-based signing and verification of env maps.
// A signature is a hex-encoded HMAC-SHA256 digest of the canonical
// (sorted key=value) representation of the env map.
package envsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Sign returns a hex-encoded HMAC-SHA256 signature for the given env map
// using the provided secret key.
func Sign(env map[string]string, secret string) string {
	canonical := canonicalize(env)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(canonical))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify returns true when the provided signature matches the env map
// signed with the given secret.
func Verify(env map[string]string, secret, signature string) bool {
	expected := Sign(env, secret)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// Fingerprint returns a short (first 12 hex chars) identifier for the env map
// useful for display purposes. It does NOT use a secret.
func Fingerprint(env map[string]string) string {
	h := sha256.Sum256([]byte(canonicalize(env)))
	return hex.EncodeToString(h[:])[:12]
}

// canonicalize builds a deterministic string from the env map.
func canonicalize(env map[string]string) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String()
}
