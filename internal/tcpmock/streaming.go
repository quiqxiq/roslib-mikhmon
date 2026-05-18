package tcpmock

import (
	"errors"
	"net"
	"sync"
)

// streamRegistry track stream yang aktif: tag → conn yang register stream itu.
// EmitToStream lookup conn lewat tag, lalu kirim !re via writeSentence.
type streamRegistry struct {
	mu     sync.Mutex
	active map[string]net.Conn
}

func newStreamRegistry() *streamRegistry {
	return &streamRegistry{active: make(map[string]net.Conn)}
}

func (r *streamRegistry) add(tag string, c net.Conn) {
	r.mu.Lock()
	r.active[tag] = c
	r.mu.Unlock()
}

func (r *streamRegistry) remove(tag string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.active[tag]; !ok {
		return false
	}
	delete(r.active, tag)
	return true
}

func (r *streamRegistry) get(tag string) (net.Conn, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.active[tag]
	return c, ok
}

func (r *streamRegistry) has(tag string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.active[tag]
	return ok
}

func (r *streamRegistry) clear() {
	r.mu.Lock()
	r.active = make(map[string]net.Conn)
	r.mu.Unlock()
}

// ErrStreamNotActive dikembalikan EmitToStream kalau tag belum ada di registry
// (belum kena handler stream, atau sudah di-cancel/close).
var ErrStreamNotActive = errors.New("tcpmock: stream tag not active")

// EmitToStream push satu sentence !re tambahan ke stream yang aktif. Tag akan
// di-prepend otomatis ke words. Aman dipanggil dari goroutine test;
// tulis ke conn diserialisasi via writeMu.
//
// Pakai ini untuk simulate event yang datang setelah subscribe (mis. user
// login muncul 50ms setelah /ip/hotspot/active/print follow).
func (s *Server) EmitToStream(tag string, words ...string) error {
	conn, ok := s.streams.get(tag)
	if !ok {
		return ErrStreamNotActive
	}
	out := append([]string{"!re", ".tag=" + tag}, words...)
	return s.writeSentence(conn, out)
}

// emitReplyWithTag kirim reply ke client (conn spesifik); kalau tag != "" dan
// reply belum punya `.tag=...`, prepend secara otomatis.
func (s *Server) emitReplyWithTag(conn net.Conn, tag string, reply []string) error {
	if tag == "" {
		return s.writeSentence(conn, reply)
	}
	// Kalau reply sudah ada `.tag=...`, jangan duplikat.
	for _, w := range reply {
		if len(w) >= 5 && w[:5] == ".tag=" {
			return s.writeSentence(conn, reply)
		}
	}
	out := make([]string, 0, len(reply)+1)
	if len(reply) > 0 {
		out = append(out, reply[0])
		out = append(out, ".tag="+tag)
		out = append(out, reply[1:]...)
	} else {
		out = append(out, ".tag="+tag)
	}
	return s.writeSentence(conn, out)
}
