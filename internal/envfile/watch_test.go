package envfile

import (
	"os"
	"testing"
	"time"
)

func writeTempEnvForWatch(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestWatchDetectsChange(t *testing.T) {
	path := writeTempEnvForWatch(t, "FOO=bar\n")
	done := make(chan struct{})
	defer close(done)

	ch, err := Watch(path, 20*time.Millisecond, done)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("FOO=baz\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-ch:
		if ev.Path != path {
			t.Errorf("expected path %s, got %s", path, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for watch event")
	}
}

func TestWatchNoEventWhenUnchanged(t *testing.T) {
	path := writeTempEnvForWatch(t, "FOO=bar\n")
	done := make(chan struct{})
	defer close(done)

	ch, err := Watch(path, 20*time.Millisecond, done)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-ch:
		t.Errorf("unexpected event: %+v", ev)
	case <-time.After(120 * time.Millisecond):
		// expected: no change detected
	}
}

func TestWatchMissingFile(t *testing.T) {
	_, err := Watch("/nonexistent/.env", 20*time.Millisecond, make(chan struct{}))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWatchStopsOnDone(t *testing.T) {
	path := writeTempEnvForWatch(t, "FOO=bar\n")
	done := make(chan struct{})

	ch, err := Watch(path, 20*time.Millisecond, done)
	if err != nil {
		t.Fatal(err)
	}

	close(done)

	// After done is closed the channel should be closed with no spurious events.
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed after done signal")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for channel to close after done")
	}
}
