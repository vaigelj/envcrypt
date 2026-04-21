package envfile

import (
	"strings"
	"testing"

	"github.com/user/envcrypt/internal/crypto"
)

func makeKey(t *testing.T) []byte {
	t.Helper()
	k, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	return k
}

func TestEncryptFieldsAll(t *testing.T) {
	key := makeKey(t)
	entries := []Entry{{Key: "DB_PASS", Value: "secret"}, {Key: "API_KEY", Value: "abc123"}}
	enc, err := EncryptFields(entries, key)
	if err != nil {
		t.Fatalf("EncryptFields: %v", err)
	}
	for _, e := range enc {
		if !strings.HasPrefix(e.Value, encryptedPrefix) {
			t.Errorf("expected encrypted prefix for %s, got %s", e.Key, e.Value)
		}
	}
}

func TestEncryptFieldsSubset(t *testing.T) {
	key := makeKey(t)
	entries := []Entry{{Key: "DB_PASS", Value: "secret"}, {Key: "HOST", Value: "localhost"}}
	enc, err := EncryptFields(entries, key, "DB_PASS")
	if err != nil {
		t.Fatalf("EncryptFields: %v", err)
	}
	if !IsEncrypted(enc[0].Value) {
		t.Errorf("DB_PASS should be encrypted")
	}
	if IsEncrypted(enc[1].Value) {
		t.Errorf("HOST should NOT be encrypted")
	}
}

func TestDecryptFieldsRoundtrip(t *testing.T) {
	key := makeKey(t)
	original := []Entry{{Key: "SECRET", Value: "topsecret"}, {Key: "PLAIN", Value: "visible"}}
	enc, err := EncryptFields(original, key, "SECRET")
	if err != nil {
		t.Fatalf("EncryptFields: %v", err)
	}
	dec, err := DecryptFields(enc, key)
	if err != nil {
		t.Fatalf("DecryptFields: %v", err)
	}
	for i, e := range dec {
		if e.Value != original[i].Value {
			t.Errorf("key %s: want %q got %q", e.Key, original[i].Value, e.Value)
		}
	}
}

func TestEncryptAlreadyEncryptedIsIdempotent(t *testing.T) {
	key := makeKey(t)
	entries := []Entry{{Key: "X", Value: "val"}}
	enc1, _ := EncryptFields(entries, key)
	enc2, err := EncryptFields(enc1, key)
	if err != nil {
		t.Fatalf("second EncryptFields: %v", err)
	}
	if enc1[0].Value != enc2[0].Value {
		t.Errorf("re-encrypting changed the value")
	}
}

func TestEncryptedKeys(t *testing.T) {
	key := makeKey(t)
	entries := []Entry{{Key: "A", Value: "v1"}, {Key: "B", Value: "v2"}, {Key: "C", Value: "v3"}}
	enc, _ := EncryptFields(entries, key, "A", "C")
	keys := EncryptedKeys(enc)
	if len(keys) != 2 {
		t.Errorf("want 2 encrypted keys, got %d", len(keys))
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := makeKey(t)
	key2 := makeKey(t)
	entries := []Entry{{Key: "S", Value: "secret"}}
	enc, _ := EncryptFields(entries, key1)
	_, err := DecryptFields(enc, key2)
	if err == nil {
		t.Error("expected error decrypting with wrong key")
	}
}
