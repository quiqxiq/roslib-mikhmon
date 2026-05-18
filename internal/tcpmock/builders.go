package tcpmock

import (
	"strings"
	"time"
)

// DoneReply mengembalikan sentence `!done` opsional dengan extra kv pairs.
// extras adalah word verbatim, mis. "=ret=*42" atau "=message=ok".
func DoneReply(extras ...string) []string {
	out := make([]string, 0, 1+len(extras))
	out = append(out, "!done")
	out = append(out, extras...)
	return out
}

// ReReply mengembalikan sentence `!re` dengan kv pairs. kvs harus berbentuk
// word lengkap (mis. "=name=foo"). Tag (.tag=...) di-prepend otomatis oleh
// streaming layer; jangan sertakan manual.
func ReReply(kvs ...string) []string {
	out := make([]string, 0, 1+len(kvs))
	out = append(out, "!re")
	out = append(out, kvs...)
	return out
}

// TrapReply mengembalikan sentence `!trap` standar dengan message+category.
// Kategori RouterOS umum: "0" (missing arg), "7" (script error). Kalau
// category kosong, tidak di-include.
func TrapReply(message, category string) []string {
	out := []string{"!trap", "=message=" + message}
	if category != "" {
		out = append(out, "=category="+category)
	}
	return out
}

// FatalReply mengembalikan sentence `!fatal` (koneksi level error).
func FatalReply(message string) []string {
	return []string{"!fatal", "=message=" + message}
}

// ActiveLoggedIn shortcut untuk row aktif hotspot dengan field lengkap.
// id dipakai sebagai `.id` (mis. "*1"), tidak perlu prefix `*`.
func ActiveLoggedIn(id, user, addr, mac string) []string {
	return []string{
		"!re",
		"=.id=" + ensureStar(id),
		"=user=" + user,
		"=address=" + addr,
		"=mac-address=" + mac,
	}
}

// ActiveLoggedOut shortcut untuk row aktif yang di-disconnect. RouterOS API
// mengirim `=.dead=true` untuk delta logout di follow mode.
func ActiveLoggedOut(id string) []string {
	return []string{
		"!re",
		"=.id=" + ensureStar(id),
		"=.dead=true",
	}
}

// UserPrint shortcut untuk row user hotspot.
func UserPrint(id, name, profile, comment string) []string {
	return []string{
		"!re",
		"=.id=" + ensureStar(id),
		"=name=" + name,
		"=profile=" + profile,
		"=comment=" + comment,
	}
}

// UserPrintExpired sama dengan UserPrint tapi comment di-set ke timestamp
// format Mikhmon (`jan/02/2006 15:04:05`, lowercase) di waktu t.
func UserPrintExpired(id, name, profile string, t time.Time) []string {
	comment := strings.ToLower(t.Format("Jan/02/2006 15:04:05"))
	return UserPrint(id, name, profile, comment)
}

// ensureStar memastikan id punya prefix `*` ala RouterOS `.id`. Idempotent.
func ensureStar(id string) string {
	if strings.HasPrefix(id, "*") {
		return id
	}
	return "*" + id
}
