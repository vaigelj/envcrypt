// Package keystore provides a simple persistent store for named AES-256
// encryption keys used by envcrypt.
//
// Keys are stored as hex-encoded strings in a JSON file on disk.
// The file is written with 0600 permissions to prevent unauthorized access.
//
// Example usage:
//
//	ks, err := keystore.New(".envcrypt/keys.json")
//	if err != nil { ... }
//
//	key, _ := crypto.GenerateKey()
//	_ = ks.Set("production", key)
//
//	retrieved, _ := ks.Get("production")
package keystore
