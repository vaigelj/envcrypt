// Package keystore manages encryption keys for envcrypt,
// supporting storage and retrieval of named keys.
package keystore

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// KeyStore holds named encryption keys persisted to disk.
type KeyStore struct {
	path string
	Keys map[string]string `json:"keys"`
}

// New loads or creates a KeyStore at the given file path.
func New(path string) (*KeyStore, error) {
	ks := &KeyStore{path: path, Keys: make(map[string]string)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return ks, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, ks); err != nil {
		return nil, err
	}
	return ks, nil
}

// Set stores a raw 32-byte key under the given name.
func (ks *KeyStore) Set(name string, key []byte) error {
	if len(key) != 32 {
		return errors.New("key must be 32 bytes")
	}
	ks.Keys[name] = hex.EncodeToString(key)
	return ks.save()
}

// Get retrieves the raw key bytes for the given name.
func (ks *KeyStore) Get(name string) ([]byte, error) {
	hex_str, ok := ks.Keys[name]
	if !ok {
		return nil, errors.New("key not found: " + name)
	}
	return hex.DecodeString(hex_str)
}

// Delete removes a key by name.
func (ks *KeyStore) Delete(name string) error {
	if _, ok := ks.Keys[name]; !ok {
		return errors.New("key not found: " + name)
	}
	delete(ks.Keys, name)
	return ks.save()
}

// List returns all stored key names.
func (ks *KeyStore) List() []string {
	names := make([]string, 0, len(ks.Keys))
	for k := range ks.Keys {
		names = append(names, k)
	}
	return names
}

func (ks *KeyStore) save() error {
	if err := os.MkdirAll(filepath.Dir(ks.path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ks.path, data, 0600)
}
