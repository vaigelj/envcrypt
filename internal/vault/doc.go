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
//
// # Error Handling
//
// All operations return descriptive errors that include the operation name and
// relevant context (e.g., key name or file path). Callers should check for
// [ErrKeyNotFound] when a named key does not exist in the keystore, and
// [ErrInvalidCiphertext] when decryption fails due to corrupt or tampered data.
package vault
