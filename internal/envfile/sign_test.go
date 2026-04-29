package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func makeSomeEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "db.example.com"},
		{Key: "SECRET", Value: "s3cr3t"},
	}
}

func TestSignAndVerifyRoundtrip(t *testing.T) {
	key := []byte("test-hmac-key-32-bytes-long!!!!")
	entries := makeSomeEntries()
	env, err := Sign(entries, key)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if env.Signature == "" {
		t.Fatal("expected non-empty signature")
	}
	if err := Verify(env, key); err != nil {
		t.Fatalf("Verify: %v", err)
	}
}

func TestVerifyFailsOnTamperedEntries(t *testing.T) {
	key := []byte("test-hmac-key-32-bytes-long!!!!")
	env, _ := Sign(makeSomeEntries(), key)
	env.Entries[0].Value = "tampered"
	if err := Verify(env, key); err == nil {
		t.Fatal("expected verification to fail after tampering")
	}
}

func TestVerifyFailsOnWrongKey(t *testing.T) {
	env, _ := Sign(makeSomeEntries(), []byte("correct-key"))
	if err := Verify(env, []byte("wrong-key")); err == nil {
		t.Fatal("expected verification to fail with wrong key")
	}
}

func TestSignEmptyKeyError(t *testing.T) {
	_, err := Sign(makeSomeEntries(), nil)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestSignFileAndVerifyFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	outPath := filepath.Join(dir, ".env.sig")
	_ = os.WriteFile(envPath, []byte("APP_ENV=production\nDB_HOST=localhost\n"), 0o644)

	key := []byte("file-hmac-key-32-bytes-long!!!!!")
	if err := SignFile(envPath, outPath, key); err != nil {
		t.Fatalf("SignFile: %v", err)
	}

	env, err := VerifyFile(outPath, key)
	if err != nil {
		t.Fatalf("VerifyFile: %v", err)
	}
	if len(env.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(env.Entries))
	}
}

func TestVerifyFileDetectsTampering(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	outPath := filepath.Join(dir, ".env.sig")
	_ = os.WriteFile(envPath, []byte("KEY=value\n"), 0o644)
	key := []byte("tamper-test-key")
	_ = SignFile(envPath, outPath, key)

	// Tamper with the envelope on disk.
	data, _ := os.ReadFile(outPath)
	var env SignedEnvelope
	_ = json.Unmarshal(data, &env)
	env.Entries = append(env.Entries, Entry{Key: "INJECTED", Value: "bad"})
	tampered, _ := json.Marshal(env)
	_ = os.WriteFile(outPath, tampered, 0o600)

	if _, err := VerifyFile(outPath, key); err == nil {
		t.Fatal("expected VerifyFile to detect tampering")
	}
}
