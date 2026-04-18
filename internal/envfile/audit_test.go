package envfile

import (
	"testing"
	"time"
)

func TestRecordAndLast(t *testing.T) {
	var log AuditLog
	if log.Last() != nil {
		t.Fatal("expected nil for empty log")
	}

	log.Record("encrypt", ".env", "alice", "ok", 5)
	e := log.Last()
	if e == nil {
		t.Fatal("expected event, got nil")
	}
	if e.Operation != "encrypt" {
		t.Errorf("expected operation encrypt, got %s", e.Operation)
	}
	if e.KeyCount != 5 {
		t.Errorf("expected 5 keys, got %d", e.KeyCount)
	}
	if e.User != "alice" {
		t.Errorf("expected user alice, got %s", e.User)
	}
}

func TestRecordMultiple(t *testing.T) {
	var log AuditLog
	log.Record("encrypt", ".env", "alice", "", 3)
	log.Record("decrypt", ".env.enc", "bob", "", 3)
	log.Record("rotate", ".env.enc", "alice", "key rotated", 3)

	if len(log.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(log.Events))
	}
	if log.Last().Operation != "rotate" {
		t.Errorf("expected last op to be rotate")
	}
}

func TestFilterByOperation(t *testing.T) {
	var log AuditLog
	log.Record("encrypt", ".env", "alice", "", 2)
	log.Record("decrypt", ".env.enc", "bob", "", 2)
	log.Record("encrypt", ".env.staging", "alice", "", 4)

	enc := log.FilterByOperation("encrypt")
	if len(enc) != 2 {
		t.Errorf("expected 2 encrypt events, got %d", len(enc))
	}
	dec := log.FilterByOperation("decrypt")
	if len(dec) != 1 {
		t.Errorf("expected 1 decrypt event, got %d", len(dec))
	}
	none := log.FilterByOperation("rotate")
	if len(none) != 0 {
		t.Errorf("expected 0 rotate events, got %d", len(none))
	}
}

func TestAuditEventString(t *testing.T) {
	e := AuditEvent{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Operation: "encrypt",
		File:      ".env",
		KeyCount:  7,
		User:      "dev",
		Details:   "success",
	}
	s := e.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
	for _, sub := range []string{"encrypt", ".env", "dev", "success", "2024"} {
		if !containsStr(s, sub) {
			t.Errorf("expected string to contain %q, got: %s", sub, s)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
