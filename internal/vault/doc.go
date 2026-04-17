// Package vault provides high-level operations for encrypting and decrypting
// .env files using the crypto and keystore packages.
//
// A Vault is tied to a keystore file that holds named AES-256 keys. Typical
// usage:
//
//	v, err := vault.New("/path/to/keys.json")
//	encrypted, err := v.EncryptFile(".env", "prod-key")
//	decrypted, err := v.DecryptFile(encrypted, "prod-key")
//
// Key rotation re-encrypts all values from one key to another without
// exposing plaintext to disk:
//
//	rotated, err := v.RotateKey(encrypted, "old-key", "new-key")
package vault
