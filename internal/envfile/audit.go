package envfile

import (
	"fmt"
	"time"
)

// AuditEvent represents a single audit log entry for an env file operation.
type AuditEvent struct {
	Timestamp time.Time
	Operation string
	File      string
	KeyCount  int
	User      string
	Details   string
}

// String returns a human-readable representation of the audit event.
func (e AuditEvent) String() string {
	return fmt.Sprintf("[%s] op=%s file=%s keys=%d user=%s details=%s",
		e.Timestamp.Format(time.RFC3339),
		e.Operation,
		e.File,
		e.KeyCount,
		e.User,
		e.Details,
	)
}

// AuditLog holds a sequence of audit events.
type AuditLog struct {
	Events []AuditEvent
}

// Record appends a new event to the audit log.
func (a *AuditLog) Record(op, file, user, details string, keyCount int) {
	a.Events = append(a.Events, AuditEvent{
		Timestamp: time.Now().UTC(),
		Operation: op,
		File:      file,
		KeyCount:  keyCount,
		User:      user,
		Details:   details,
	})
}

// FilterByOperation returns all events matching the given operation name.
func (a *AuditLog) FilterByOperation(op string) []AuditEvent {
	var out []AuditEvent
	for _, e := range a.Events {
		if e.Operation == op {
			out = append(out, e)
		}
	}
	return out
}

// Last returns the most recent audit event, or nil if the log is empty.
func (a *AuditLog) Last() *AuditEvent {
	if len(a.Events) == 0 {
		return nil
	}
	e := a.Events[len(a.Events)-1]
	return &e
}
