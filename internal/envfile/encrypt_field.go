package envfile

import (
	"errors"
	"fmt"
	"strings"

	"github.com/user/envcrypt/internal/crypto"
)

const encryptedPrefix = "enc:"

// EncryptFields encrypts the values of the given keys in entries using the
// provided 32-byte key. Encrypted values are stored with the "enc:" prefix.
func EncryptFields(entries []Entry, key []byte, keys ...string) ([]Entry, error) {
	keySet := toSet(keys)
	result := make([]Entry, len(entries))
	for i, e := range entries {
		if len(keySet) > 0 && !keySet[e.Key] {
			result[i] = e
			continue
		}
		if strings.HasPrefix(e.Value, encryptedPrefix) {
			// already encrypted
			result[i] = e
			continue
		}
		ciphertext, err := crypto.Encrypt(key, []byte(e.Value))
		if err != nil {
			return nil, fmt.Errorf("encrypt field %q: %w", e.Key, err)
		}
		result[i] = Entry{Key: e.Key, Value: encryptedPrefix + ciphertext}
	}
	return result, nil
}

// DecryptFields decrypts values prefixed with "enc:" in entries using the
// provided 32-byte key.
func DecryptFields(entries []Entry, key []byte) ([]Entry, error) {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		if !strings.HasPrefix(e.Value, encryptedPrefix) {
			result[i] = e
			continue
		}
		encoded := strings.TrimPrefix(e.Value, encryptedPrefix)
		plaintext, err := crypto.Decrypt(key, encoded)
		if err != nil {
			return nil, fmt.Errorf("decrypt field %q: %w", e.Key, err)
		}
		result[i] = Entry{Key: e.Key, Value: string(plaintext)}
	}
	return result, nil
}

// IsEncrypted reports whether the given entry value is encrypted.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}

// EncryptedKeys returns the keys whose values are currently encrypted.
func EncryptedKeys(entries []Entry) []string {
	var out []string
	for _, e := range entries {
		if IsEncrypted(e.Value) {
			out = append(out, e.Key)
		}
	}
	return out
}

// ErrNoEncryptedFields is returned when no encrypted fields are found.
var ErrNoEncryptedFields = errors.New("no encrypted fields found")
