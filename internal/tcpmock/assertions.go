package tcpmock

import (
	"testing"
	"time"
)

// AssertReceived memverifikasi setidaknya satu sentence yang sudah masuk
// cocok dengan matcher. Mengembalikan sentence yang match. Polling sampai 1s
// supaya non-flaky terhadap latency goroutine accept.
func (s *Server) AssertReceived(t *testing.T, m Matcher, msg string) []string {
	t.Helper()
	deadline := time.Now().Add(1 * time.Second)
	for time.Now().Before(deadline) {
		for _, w := range s.Received() {
			if m(w) {
				return w
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("tcpmock: expected sentence not received: %s\nrecorded: %v", msg, s.Received())
	return nil
}

// AssertReceivedAll memverifikasi tepat count sentence yang match. Polling
// sampai 1s.
func (s *Server) AssertReceivedAll(t *testing.T, m Matcher, count int) [][]string {
	t.Helper()
	deadline := time.Now().Add(1 * time.Second)
	var matched [][]string
	for time.Now().Before(deadline) {
		matched = matched[:0]
		for _, w := range s.Received() {
			if m(w) {
				matched = append(matched, w)
			}
		}
		if len(matched) >= count {
			return matched[:count]
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("tcpmock: expected %d matching sentences, got %d\nrecorded: %v", count, len(matched), s.Received())
	return nil
}

// AssertNotReceived memverifikasi tidak ada sentence yang match. Tunggu 200ms
// sebagai jaminan tidak ada race "belum sempat datang".
func (s *Server) AssertNotReceived(t *testing.T, m Matcher) {
	t.Helper()
	time.Sleep(200 * time.Millisecond)
	for _, w := range s.Received() {
		if m(w) {
			t.Fatalf("tcpmock: did not expect sentence but got: %v", w)
		}
	}
}
