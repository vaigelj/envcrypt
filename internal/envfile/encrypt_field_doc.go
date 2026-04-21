// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// # Field-level Encryption
//
// EncryptFields and DecryptFields allow individual values within an env file
// to be encrypted using AES-256-GCM (via the crypto package). Encrypted
// values are stored inline with an "enc:" prefix so that the file remains
// valid and partially human-readable.
//
// Example workflow:
//
//	entries, _ := Parse(r)
//	key, _     := crypto.GenerateKey()
//	enc, _     := EncryptFields(entries, key, "DB_PASSWORD", "API_SECRET")
//	// write enc to disk; only those two values are opaque
//	dec, _     := DecryptFields(enc, key)
package envfile
