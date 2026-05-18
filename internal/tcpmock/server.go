package tcpmock

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

// Server adalah TCP server minimal RouterOS API.
//
// Mode dispatch:
//   - Handler matcher-based via OnSentence/OnStream (cek dulu).
//   - FIFO scripted via Script (fallback kalau tidak ada handler match).
//
// Multi-connection: tiap accepted conn dilayani di goroutine sendiri. Cocok
// untuk roslib.Device yang mendial dua koneksi (stream + command).
type Server struct {
	ln net.Listener

	mu        sync.Mutex
	received  [][]string
	scripted  [][][]string // FIFO fallback grup; satu grup = semua reply untuk satu command
	handlers  []handler
	closeOnce sync.Once
	conns     []net.Conn
	done      chan struct{}

	// writeMu menyerialisasi semua write ke conn supaya EmitToStream (test
	// goroutine) tidak race dengan reply dari accept loop.
	writeMu sync.Mutex

	streams *streamRegistry
}

// Start membuka listener pada 127.0.0.1:0 dan kembalikan Server siap pakai.
// Goroutine accept dimulai langsung dan menerima conn tanpa batas.
func Start() (*Server, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("tcpmock: listen: %w", err)
	}
	s := &Server{
		ln:      ln,
		done:    make(chan struct{}),
		streams: newStreamRegistry(),
	}
	go s.acceptLoop()
	return s, nil
}

// Addr alamat listener (host:port).
func (s *Server) Addr() string { return s.ln.Addr().String() }

// Script men-queue satu grup reply yang akan dikirim ke client sebagai
// fallback kalau tidak ada handler yang match. Satu Script call = satu grup
// = semua reply untuk SATU command (mis. !re+!re+!done untuk Print). Boleh
// variadic: `Script(re1, re2, done)` mengirim tiga sentence berurutan saat
// sentence non-handler berikutnya masuk.
func (s *Server) Script(replies ...[]string) {
	if len(replies) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scripted = append(s.scripted, replies)
}

// Received mengembalikan snapshot semua sentence yang sudah masuk dari semua
// conn (gabungan, urutan masuk).
func (s *Server) Received() [][]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([][]string, len(s.received))
	for i, r := range s.received {
		out[i] = append([]string(nil), r...)
	}
	return out
}

// Close menutup koneksi dan listener. Aman dipanggil berkali-kali.
func (s *Server) Close() error {
	s.closeOnce.Do(func() {
		_ = s.ln.Close()
		s.mu.Lock()
		conns := append([]net.Conn(nil), s.conns...)
		s.conns = nil
		s.mu.Unlock()
		for _, c := range conns {
			_ = c.Close()
		}
		if s.streams != nil {
			s.streams.clear()
		}
		close(s.done)
	})
	return nil
}

// writeSentence kirim sentence ke conn target dengan locking. Boleh dipanggil
// dari accept goroutine atau dari test goroutine (EmitToStream).
func (s *Server) writeSentence(c net.Conn, words []string) error {
	if c == nil {
		return errors.New("tcpmock: nil conn")
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return EncodeSentence(c, words)
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		s.mu.Lock()
		s.conns = append(s.conns, conn)
		s.mu.Unlock()
		go s.serve(conn)
	}
}

func (s *Server) serve(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		words, err := DecodeSentence(r)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			return
		}

		s.mu.Lock()
		s.received = append(s.received, words)
		s.mu.Unlock()

		s.dispatch(conn, words)
	}
}

// dispatch jalankan handler yang cocok atau FIFO fallback. Tag echo +
// streaming registry di-handle di sini supaya semua write lewat satu jalur.
func (s *Server) dispatch(conn net.Conn, words []string) {
	tag := extractTag(words)

	// Handle /cancel: remove tag dari registry, reply !done dengan tag yang
	// di-cancel di-echo (RouterOS pattern).
	if len(words) > 0 && words[0] == "/cancel" {
		targetTag := extractEqArg(words, "tag")
		if targetTag != "" {
			s.streams.remove(targetTag)
			// kirim !done untuk command /cancel sendiri (pakai its own tag)
			_ = s.emitReplyWithTag(conn, tag, []string{"!done"})
			// kirim !done dengan tag asli supaya listener selesai (RouterOS pattern)
			_ = s.writeSentence(conn, []string{"!done", ".tag=" + targetTag})
		}
		return
	}

	s.mu.Lock()
	handlers := s.handlers
	s.mu.Unlock()

	for _, h := range handlers {
		if !h.match(words) {
			continue
		}
		if h.isStream && tag == "" {
			// Stream handler match tapi sentence tidak punya .tag — skip.
			continue
		}
		replies := h.replies
		if h.fn != nil {
			replies = h.fn(words)
		}
		for _, reply := range replies {
			_ = s.emitReplyWithTag(conn, tag, reply)
		}
		if h.isStream {
			s.streams.add(tag, conn)
		}
		return
	}

	// Fallback: FIFO scripted (grup berisi N sentences untuk satu command).
	s.mu.Lock()
	var group [][]string
	if len(s.scripted) > 0 {
		group = s.scripted[0]
		s.scripted = s.scripted[1:]
	}
	s.mu.Unlock()
	for _, reply := range group {
		_ = s.emitReplyWithTag(conn, tag, reply)
	}
}
