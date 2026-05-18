package testutil

import (
	"testing"
	"time"
)

// Eventually polling sampai cond() true atau timeout. Fatal kalau timeout.
// tickInterval 10ms — cukup untuk test fan-out goroutine tanpa busy-loop.
func Eventually(t *testing.T, cond func() bool, timeout time.Duration, msg string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("eventually: %s (timeout %v)", msg, timeout)
}
