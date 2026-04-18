package envfile

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent is emitted when a watched file changes.
type WatchEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Watch polls the given file at the specified interval and sends a WatchEvent
// on the returned channel whenever the file content changes. The caller must
// close done to stop watching.
func Watch(path string, interval time.Duration, done <-chan struct{}) (<-chan WatchEvent, error) {
	initialHash, err := fileHash(path)
	if err != nil {
		return nil, fmt.Errorf("watch: %w", err)
	}

	ch := make(chan WatchEvent, 1)
	go func() {
		defer close(ch)
		last := initialHash
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				current, err := fileHash(path)
				if err != nil {
					continue
				}
				if current != last {
					ch <- WatchEvent{Path: path, OldHash: last, NewHash: current, At: t}
					last = current
				}
			}
		}
	}()
	return ch, nil
}
