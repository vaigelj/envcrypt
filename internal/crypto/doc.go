// Package crypto provides AES-256-GCM encryption and decryption primitives
// used by envcrypt to protect .env file contents.
//
// Usage:
//
//	key, err := crypto.GenerateKey()
//	if err != nil { ... }
//
//	ciphertext, err := crypto.Encrypt(key, []byte("MY_SECRET=value"))
//	if err != nil { ... }
//
//	plaintext, err := crypto.Decrypt(key, ciphertext)
//	if err != nil { ... }
//
// Keys are 32 bytes (256 bits). Nonces are randomly generated per encryption
// and prepended to the ciphertext output, so each call to Encrypt produces
// a unique ciphertext even for identical inputs.
package crypto
