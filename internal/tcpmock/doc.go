// Package tcpmock menyediakan server TCP minimal yang berbicara protocol
// RouterOS API (length-prefix word encoding) untuk test encoding-level dan
// emulate sentence flow tanpa router fisik.
//
// Fitur:
//
//   - Encoding/decoding word & sentence sesuai spec RouterOS API.
//   - Matcher-based dispatch (OnSentence): registrasi handler yang reply
//     berdasarkan pattern sentence (command path, kv arg, dll).
//   - Streaming (OnStream + EmitToStream): handler yang track `.tag=X` dan
//     pertahankan stream sampai /cancel atau Close. Test bisa push !re
//     event tambahan via EmitToStream untuk simulate event yang datang
//     setelah subscribe.
//   - Login (AcceptLogin): handler default untuk /login modern (post-6.43).
//   - Sentence builder helpers (DoneReply, ReReply, ActiveLoggedIn, dst)
//     untuk test setup yang readable.
//   - Assertion helpers (AssertReceived, AssertReceivedAll, AssertNotReceived).
//   - FIFO Script (fallback) untuk test legacy yang belum migrate.
//
// Batasan sengaja:
//
//   - Single koneksi (satu test = satu device).
//   - Tidak ada legacy login challenge/response (pre-6.43).
//
// Untuk e2e ke router fisik pakai paket `test/integration/`.
package tcpmock
