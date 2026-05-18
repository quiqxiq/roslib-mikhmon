package tcpmock

// AcceptLogin mendaftarkan handler default untuk `/login` ala RouterOS post-6.43
// (modern flow tanpa challenge-response). Sentence `/login =name=X =password=Y`
// langsung di-reply `!done` tanpa validasi credential.
//
// Legacy login flow (pre-6.43 challenge/response) tidak didukung — test
// infrastructure mengasumsikan auth modern. Kalau butuh, daftarkan handler
// manual via OnSentence.
func (s *Server) AcceptLogin() *Server {
	return s.OnSentence(MatchCommand("/login"), DoneReply())
}
