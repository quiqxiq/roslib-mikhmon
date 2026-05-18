package tcpmock

import "strings"

// Matcher mengembalikan true kalau sentence cocok dengan pola handler.
type Matcher func(words []string) bool

// handler adalah satu entri tabel dispatch. Kalau Match cocok dengan sentence
// yang masuk, server akan mengirim Replies (atau panggil fn kalau diset).
type handler struct {
	match    Matcher
	replies  [][]string
	fn       func(words []string) [][]string // optional dynamic replies; menang atas replies
	isStream bool
}

// OnSentence mendaftarkan handler immediate-reply: setiap kali sentence masuk
// dan cocok dengan matcher, server mengirim replies berurutan ke client. Tag
// (`.tag=...`) di sentence di-echo otomatis ke tiap reply.
//
// Handler dievaluasi sesuai urutan registrasi; handler pertama yang match
// menang. Kalau tidak ada handler match, server fallback ke FIFO Script.
func (s *Server) OnSentence(m Matcher, replies ...[]string) *Server {
	s.mu.Lock()
	s.handlers = append(s.handlers, handler{match: m, replies: replies})
	s.mu.Unlock()
	return s
}

// OnStream mendaftarkan handler streaming: sentence yang cocok harus punya
// `.tag=X`. Server mengirim taggedReplies (tag di-prepend) dan mempertahankan
// tag itu di registry, sehingga EmitToStream dapat push !re berikutnya sampai
// client mengirim /cancel atau Server.Close().
func (s *Server) OnStream(m Matcher, taggedReplies ...[]string) *Server {
	s.mu.Lock()
	s.handlers = append(s.handlers, handler{match: m, replies: taggedReplies, isStream: true})
	s.mu.Unlock()
	return s
}

// OnSentenceFunc seperti OnSentence tapi reply ditentukan dinamis oleh fn
// berdasarkan sentence yang masuk (mis. counter, state machine). Berguna
// untuk skenario partial-failure di mana N panggilan pertama OK dan sisanya
// trap.
func (s *Server) OnSentenceFunc(m Matcher, fn func(words []string) [][]string) *Server {
	s.mu.Lock()
	s.handlers = append(s.handlers, handler{match: m, fn: fn})
	s.mu.Unlock()
	return s
}

// MatchCommand match kalau word pertama (command path) sama persis.
func MatchCommand(cmd string) Matcher {
	return func(words []string) bool {
		return len(words) > 0 && words[0] == cmd
	}
}

// MatchHas match kalau ada word berbentuk "=key=value" persis.
func MatchHas(key, value string) Matcher {
	prefix := "=" + key + "="
	return func(words []string) bool {
		for _, w := range words {
			if strings.HasPrefix(w, prefix) && w[len(prefix):] == value {
				return true
			}
		}
		return false
	}
}

// MatchHasKey match kalau ada word "=key=..." (value bebas).
func MatchHasKey(key string) Matcher {
	prefix := "=" + key + "="
	return func(words []string) bool {
		for _, w := range words {
			if strings.HasPrefix(w, prefix) {
				return true
			}
		}
		return false
	}
}

// MatchAll AND-kan beberapa matcher.
func MatchAll(ms ...Matcher) Matcher {
	return func(words []string) bool {
		for _, m := range ms {
			if !m(words) {
				return false
			}
		}
		return true
	}
}

// extractTag mengembalikan value `.tag=X` di sentence, atau "" kalau tidak ada.
func extractTag(words []string) string {
	for _, w := range words {
		if strings.HasPrefix(w, ".tag=") {
			return w[len(".tag="):]
		}
	}
	return ""
}

// extractEqArg mengembalikan value untuk word "=key=value" pertama, atau "".
func extractEqArg(words []string, key string) string {
	prefix := "=" + key + "="
	for _, w := range words {
		if strings.HasPrefix(w, prefix) {
			return w[len(prefix):]
		}
	}
	return ""
}
