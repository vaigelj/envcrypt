package envfile

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"time"
)

// SignedEnvelope wraps a set of env entries with an HMAC signature and metadata.
type SignedEnvelope struct {
	Entries   []Entry   `json:"entries"`
	SignedAt  time.Time `json:"signed_at"`
	Signature string    `json:"signature"`
}

// Sign creates a SignedEnvelope for the given entries using the provided HMAC key.
func Sign(entries []Entry, key []byte) (*SignedEnvelope, error) {
	if len(key) == 0 {
		return nil, errors.New("sign: key must not be empty")
	}
	now := time.Now().UTC()
	env := &SignedEnvelope{
		Entries:  entries,
		SignedAt: now,
	}
	sig, err := computeSignature(env.Entries, env.SignedAt, key)
	if err != nil {
		return nil, err
	}
	env.Signature = sig
	return env, nil
}

// Verify checks whether the envelope's signature is valid for the given key.
func Verify(env *SignedEnvelope, key []byte) error {
	if len(key) == 0 {
		return errors.New("verify: key must not be empty")
	}
	expected, err := computeSignature(env.Entries, env.SignedAt, key)
	if err != nil {
		return err
	}
	if !hmac.Equal([]byte(env.Signature), []byte(expected)) {
		return errors.New("verify: signature mismatch — envelope may have been tampered")
	}
	return nil
}

// SignFile reads a .env file, signs its entries, and writes the envelope as JSON.
func SignFile(envPath, outPath string, key []byte) error {
	entries, err := ParseFile(envPath)
	if err != nil {
		return err
	}
	env, err := Sign(entries, key)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, data, 0o600)
}

// VerifyFile loads a signed envelope from disk and verifies it.
func VerifyFile(envelopePath string, key []byte) (*SignedEnvelope, error) {
	data, err := os.ReadFile(envelopePath)
	if err != nil {
		return nil, err
	}
	var env SignedEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	if err := Verify(&env, key); err != nil {
		return nil, err
	}
	return &env, nil
}

func computeSignature(entries []Entry, ts time.Time, key []byte) (string, error) {
	h := hmac.New(sha256.New, key)
	payload, err := json.Marshal(struct {
		Entries  []Entry   `json:"entries"`
		SignedAt time.Time `json:"signed_at"`
	}{entries, ts})
	if err != nil {
		return "", err
	}
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil)), nil
}
