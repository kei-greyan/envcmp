// Package envsign provides utilities for signing and verifying env maps using
// HMAC-SHA256. This allows callers to detect unauthorised modifications to a
// set of environment variables between environments or pipeline stages.
//
// Usage:
//
//	sig := envsign.Sign(env, secret)
//	ok  := envsign.Verify(env, secret, sig)
//	fp  := envsign.Fingerprint(env) // secret-free short identifier
package envsign
