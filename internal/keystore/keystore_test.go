package keystore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcrypt/internal/keystore"
)

func tempStore(t *testing.T) (*keystore.KeyStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "keys.json")
	ks, err := keystore.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return ks, path
}

func TestSetAndGet(t *testing.T) {
	ks, _ := tempStore(t)
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	if err := ks.Set("default", key); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := ks.Get("default")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	for i, b := range got {
		if b != key[i] {
			t.Fatalf("key mismatch at index %d", i)
		}
	}
}

func TestPersistence(t *testing.T) {
	ks, path := tempStore(t)
	key := make([]byte, 32)
	_ = ks.Set("persist", key)

	ks2, err := keystore.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if _, err := ks2.Get("persist"); err != nil {
		t.Fatalf("key not persisted: %v", err)
	}
}

func TestDelete(t *testing.T) {
	ks, _ := tempStore(t)
	key := make([]byte, 32)
	_ = ks.Set("temp", key)
	if err := ks.Delete("temp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := ks.Get("temp"); err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestSetInvalidKeySize(t *testing.T) {
	ks, _ := tempStore(t)
	if err := ks.Set("bad", []byte{1, 2, 3}); err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestGetMissing(t *testing.T) {
	ks, _ := tempStore(t)
	if _, err := ks.Get("missing"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestFilePermissions(t *testing.T) {
	ks, path := tempStore(t)
	_ = ks.Set("perm", make([]byte, 32))
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
