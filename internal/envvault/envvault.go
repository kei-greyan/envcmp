// Package envvault provides encryption and decryption of .env file values
// using AES-GCM with a passphrase-derived key.
package envvault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const prefix = "vault:"

// Encrypt encrypts each value in env using the given passphrase.
// Already-encrypted values (prefixed with "vault:") are left unchanged.
// Returns a new map; the original is not mutated.
func Encrypt(env map[string]string, passphrase string) (map[string]string, error) {
	if passphrase == "" {
		return nil, errors.New("envvault: passphrase must not be empty")
	}
	key := deriveKey(passphrase)
	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(v) >= len(prefix) && v[:len(prefix)] == prefix {
			out[k] = v
			continue
		}
		enc, err := encryptValue(key, v)
		if err != nil {
			return nil, fmt.Errorf("envvault: encrypt key %q: %w", k, err)
		}
		out[k] = prefix + enc
	}
	return out, nil
}

// Decrypt decrypts each vault-prefixed value in env using the given passphrase.
// Values without the "vault:" prefix are passed through unchanged.
// Returns a new map; the original is not mutated.
func Decrypt(env map[string]string, passphrase string) (map[string]string, error) {
	if passphrase == "" {
		return nil, errors.New("envvault: passphrase must not be empty")
	}
	key := deriveKey(passphrase)
	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(v) < len(prefix) || v[:len(prefix)] != prefix {
			out[k] = v
			continue
		}
		dec, err := decryptValue(key, v[len(prefix):])
		if err != nil {
			return nil, fmt.Errorf("envvault: decrypt key %q: %w", k, err)
		}
		out[k] = dec
	}
	return out, nil
}

// IsEncrypted reports whether all values in env carry the vault prefix.
func IsEncrypted(env map[string]string) bool {
	for _, v := range env {
		if len(v) < len(prefix) || v[:len(prefix)] != prefix {
			return false
		}
	}
	return true
}

func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

func encryptValue(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptValue(key []byte, encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("ciphertext too short")
	}
	plaintext, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	return string(plaintext), nil
}
