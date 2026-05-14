// Package envvault provides symmetric encryption and decryption of environment
// variable values using AES-256-GCM.
//
// Encrypted values are stored with a "vault:" prefix followed by a
// base64-encoded ciphertext, making them safe to commit to version control
// without exposing secrets.
//
// The encryption key is derived from a user-supplied passphrase via SHA-256.
// For production use, consider supplying a strong, randomly generated
// passphrase stored in a secrets manager.
//
// Example usage:
//
//	enc, err := envvault.Encrypt(env, passphrase)
//	dec, err := envvault.Decrypt(enc, passphrase)
package envvault
