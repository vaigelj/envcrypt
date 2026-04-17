package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcrypt/internal/crypto"
	"github.com/user/envcrypt/internal/envfile"
	"github.com/user/envcrypt/internal/keystore"
	"github.com/user/envcrypt/internal/vault"
)

func setupVault(t *testing.T) (*vault.Vault, string) {
	t.Helper()
	dir := t.TempDir()
	ksPath := filepath.Join(dir, "keys.json")
	v, err := vault.New(ksPath)
	if err != nil {
		t.Fatalf("vault.New: %v", err)
	}
	return v, ksPath
}

func seedKey(t *testing.T, ksPath, name string) {
	t.Helper()
	s, _ := keystore.New(ksPath)
	key, _ := crypto.GenerateKey()
	if err := s.Set(name, key); err != nil {
		t.Fatalf("seed key: %v", err)
	}
}

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	v, ksPath := setupVault(t)
	seedKey(t, ksPath, "k1")

	src := writeTempEnv(t, "DB_HOST=localhost\nDB_PASS=secret\n")
	enc, err := v.EncryptFile(src, "k1")
	if err != nil {
		t.Fatalf("EncryptFile: %v", err)
	}

	dec, err := v.DecryptFile(enc, "k1")
	if err != nil {
		t.Fatalf("DecryptFile: %v", err)
	}

	if dec["DB_HOST"] != "localhost" || dec["DB_PASS"] != "secret" {
		t.Errorf("unexpected decrypted values: %v", dec)
	}
}

func TestRotateKey(t *testing.T) {
	v, ksPath := setupVault(t)
	seedKey(t, ksPath, "old")
	seedKey(t, ksPath, "new")

	src := writeTempEnv(t, "API_KEY=topsecret\n")
	enc, err := v.EncryptFile(src, "old")
	if err != nil {
		t.Fatalf("EncryptFile: %v", err)
	}

	rotated, err := v.RotateKey(enc, "old", "new")
	if err != nil {
		t.Fatalf("RotateKey: %v", err)
	}

	dec, err := v.DecryptFile(rotated, "new")
	if err != nil {
		t.Fatalf("DecryptFile after rotate: %v", err)
	}
	if dec["API_KEY"] != "topsecret" {
		t.Errorf("got %q, want topsecret", dec["API_KEY"])
	}
}

func TestEncryptMissingKey(t *testing.T) {
	v, _ := setupVault(t)
	src := writeTempEnv(t, "X=1\n")
	_, err := v.EncryptFile(src, "ghost")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestDecryptWrongKey(t *testing.T) {
	v, ksPath := setupVault(t)
	seedKey(t, ksPath, "a")
	seedKey(t, ksPath, "b")

	src := writeTempEnv(t, "VAR=value\n")
	enc, _ := v.EncryptFile(src, "a")

	// Manually build an EnvMap with wrong-key ciphertext
	_, err := v.DecryptFile(enc, "b")
	if err == nil {
		t.Error("expected decryption error with wrong key")
	}
}

var _ = envfile.EnvMap(nil) // ensure import used
