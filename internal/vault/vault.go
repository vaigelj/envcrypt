package vault

import (
	"fmt"

	"github.com/user/envcrypt/internal/crypto"
	"github.com/user/envcrypt/internal/envfile"
	"github.com/user/envcrypt/internal/keystore"
)

// Vault combines envfile, keystore, and crypto to encrypt/decrypt .env files.
type Vault struct {
	store *keystore.Store
}

// New creates a new Vault backed by the given keystore path.
func New(keystorePath string) (*Vault, error) {
	s, err := keystore.New(keystorePath)
	if err != nil {
		return nil, fmt.Errorf("vault: open keystore: %w", err)
	}
	return &Vault{store: s}, nil
}

// EncryptFile reads plaintext env vars from src, encrypts each value using
// the named key, and returns an EnvMap with ciphertext values.
func (v *Vault) EncryptFile(src string, keyName string) (envfile.EnvMap, error) {
	key, err := v.store.Get(keyName)
	if err != nil {
		return nil, fmt.Errorf("vault: get key %q: %w", keyName, err)
	}

	vars, err := envfile.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("vault: parse env file: %w", err)
	}

	result := make(envfile.EnvMap, len(vars))
	for k, v := range vars {
		cipher, err := crypto.Encrypt(key, []byte(v))
		if err != nil {
			return nil, fmt.Errorf("vault: encrypt %q: %w", k, err)
		}
		result[k] = cipher
	}
	return result, nil
}

// DecryptFile takes an EnvMap of ciphertext values, decrypts each using the
// named key, and returns a plaintext EnvMap.
func (v *Vault) DecryptFile(encrypted envfile.EnvMap, keyName string) (envfile.EnvMap, error) {
	key, err := v.store.Get(keyName)
	if err != nil {
		return nil, fmt.Errorf("vault: get key %q: %w", keyName, err)
	}

	result := make(envfile.EnvMap, len(encrypted))
	for k, cipherText := range encrypted {
		plain, err := crypto.Decrypt(key, cipherText)
		if err != nil {
			return nil, fmt.Errorf("vault: decrypt %q: %w", k, err)
		}
		result[k] = string(plain)
	}
	return result, nil
}

// RotateKey re-encrypts an EnvMap from oldKey to newKey and returns the result.
func (v *Vault) RotateKey(encrypted envfile.EnvMap, oldKey, newKey string) (envfile.EnvMap, error) {
	plain, err := v.DecryptFile(encrypted, oldKey)
	if err != nil {
		return nil, fmt.Errorf("vault: rotate decrypt: %w", err)
	}
	result := make(envfile.EnvMap, len(plain))
	nk, err := v.store.Get(newKey)
	if err != nil {
		return nil, fmt.Errorf("vault: get new key %q: %w", newKey, err)
	}
	for k, val := range plain {
		cipher, err := crypto.Encrypt(nk, []byte(val))
		if err != nil {
			return nil, fmt.Errorf("vault: rotate encrypt %q: %w", k, err)
		}
		result[k] = cipher
	}
	return result, nil
}
