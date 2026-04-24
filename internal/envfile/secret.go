// Package envfile provides utilities for parsing, writing, and manipulating .env files.
package envfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/json"
)

// SecretEntry represents a named secret with metadata.
type SecretEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SecretStore holds a collection of named secrets persisted to disk.
type SecretStore struct {
	Secrets map[string]*SecretEntry `json:"secrets"`
}

func secretsPath(dir string) string {
	return filepath.Join(dir, ".envcrypt_secrets.json")
}

// LoadSecrets reads the secret store from dir. Returns an empty store if the
// file does not exist yet.
func LoadSecrets(dir string) (*SecretStore, error) {
	path := secretsPath(dir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &SecretStore{Secrets: make(map[string]*SecretEntry)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load secrets: %w", err)
	}
	var store SecretStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse secrets: %w", err)
	}
	if store.Secrets == nil {
		store.Secrets = make(map[string]*SecretEntry)
	}
	return &store, nil
}

// SaveSecrets persists the secret store to dir.
func SaveSecrets(dir string, store *SecretStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal secrets: %w", err)
	}
	path := secretsPath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create secrets dir: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// SetSecret adds or updates a secret identified by key.
func SetSecret(store *SecretStore, key, value, note string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("secret key must not be empty")
	}
	now := time.Now().UTC()
	if existing, ok := store.Secrets[key]; ok {
		existing.Value = value
		existing.Note = note
		existing.UpdatedAt = now
		return nil
	}
	store.Secrets[key] = &SecretEntry{
		Key:       key,
		Value:     value,
		Note:      note,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return nil
}

// GetSecret retrieves a secret by key. Returns an error if not found.
func GetSecret(store *SecretStore, key string) (*SecretEntry, error) {
	entry, ok := store.Secrets[key]
	if !ok {
		return nil, fmt.Errorf("secret %q not found", key)
	}
	return entry, nil
}

// DeleteSecret removes a secret by key. Returns an error if the key does not exist.
func DeleteSecret(store *SecretStore, key string) error {
	if _, ok := store.Secrets[key]; !ok {
		return fmt.Errorf("secret %q not found", key)
	}
	delete(store.Secrets, key)
	return nil
}

// ListSecretKeys returns all secret keys in sorted order.
func ListSecretKeys(store *SecretStore) []string {
	keys := make([]string, 0, len(store.Secrets))
	for k := range store.Secrets {
		keys = append(keys, k)
	}
	// stable sort
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
