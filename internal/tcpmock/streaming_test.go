package tcpmock

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dialClient buka TCP koneksi ke mock server dan kembalikan conn + reader.
// t.Cleanup tutup conn.
func dialClient(t *testing.T, srv *Server) (net.Conn, *bufio.Reader) {
	t.Helper()
	conn, err := net.DialTimeout("tcp", srv.Addr(), 2*time.Second)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return conn, bufio.NewReader(conn)
}

func TestServer_OnSentence_dispatches(t *testing.T) {
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	srv.OnSentence(MatchCommand("/ping"), DoneReply("=ret=pong"))

	conn, r := dialClient(t, srv)
	require.NoError(t, EncodeSentence(conn, []string{"/ping"}))

	got, err := DecodeSentence(r)
	require.NoError(t, err)
	assert.Equal(t, []string{"!done", "=ret=pong"}, got)
}

func TestServer_login_handshake(t *testing.T) {
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	srv.AcceptLogin()

	conn, r := dialClient(t, srv)
	require.NoError(t, EncodeSentence(conn, []string{"/login", "=name=admin", "=password=secret"}))

	got, err := DecodeSentence(r)
	require.NoError(t, err)
	assert.Equal(t, []string{"!done"}, got)
}

func TestServer_streamReply_echoesTag(t *testing.T) {
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	srv.OnStream(MatchCommand("/ip/hotspot/active/print"), ReReply("=.id=*1", "=user=alice"))

	conn, r := dialClient(t, srv)
	require.NoError(t, EncodeSentence(conn, []string{
		"/ip/hotspot/active/print",
		"=follow=",
		".tag=T1",
	}))

	got, err := DecodeSentence(r)
	require.NoError(t, err)
	require.Len(t, got, 4)
	assert.Equal(t, "!re", got[0])
	assert.Equal(t, ".tag=T1", got[1])
	assert.Equal(t, "=.id=*1", got[2])
	assert.Equal(t, "=user=alice", got[3])
}

func TestServer_emitToStream_pushesEvent(t *testing.T) {
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	// Register stream handler tanpa initial reply — kita push event manual.
	srv.OnStream(MatchCommand("/ip/hotspot/active/print"))

	conn, r := dialClient(t, srv)
	require.NoError(t, EncodeSentence(conn, []string{
		"/ip/hotspot/active/print",
		"=follow=",
		".tag=T2",
	}))

	// Eventually push event.
	require.Eventually(t, func() bool {
		return srv.streams.has("T2")
	}, 500*time.Millisecond, 10*time.Millisecond)

	require.NoError(t, srv.EmitToStream("T2", "=.id=*5", "=user=bob"))

	got, err := DecodeSentence(r)
	require.NoError(t, err)
	assert.Equal(t, []string{"!re", ".tag=T2", "=.id=*5", "=user=bob"}, got)
}

func TestServer_emitToStream_unknownTag(t *testing.T) {
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	err = srv.EmitToStream("nonexistent", "=foo=bar")
	assert.ErrorIs(t, err, ErrStreamNotActive)
}

func TestServer_cancel_stopsStream(t *testing.T) {
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	srv.OnStream(MatchCommand("/ip/hotspot/active/print"))

	conn, r := dialClient(t, srv)
	require.NoError(t, EncodeSentence(conn, []string{
		"/ip/hotspot/active/print",
		"=follow=",
		".tag=T3",
	}))
	require.Eventually(t, func() bool {
		return srv.streams.has("T3")
	}, 500*time.Millisecond, 10*time.Millisecond)

	// Client kirim /cancel =tag=T3 dengan own tag T-cancel.
	require.NoError(t, EncodeSentence(conn, []string{
		"/cancel",
		"=tag=T3",
		".tag=T-cancel",
	}))

	// Tunggu cancel diproses.
	require.Eventually(t, func() bool {
		return !srv.streams.has("T3")
	}, 500*time.Millisecond, 10*time.Millisecond)

	// EmitToStream untuk T3 sekarang error.
	err = srv.EmitToStream("T3", "=foo=bar")
	assert.ErrorIs(t, err, ErrStreamNotActive)

	// Drain reply dari /cancel: harusnya ada !done untuk cancel & !done untuk tag asli.
	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	for i := 0; i < 2; i++ {
		s, err := DecodeSentence(r)
		if err != nil {
			break
		}
		assert.NotEmpty(t, s)
	}
}

func TestServer_assertReceived_helpers(t *testing.T) {
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	srv.OnSentence(MatchCommand("/x"), DoneReply())

	conn, _ := dialClient(t, srv)
	require.NoError(t, EncodeSentence(conn, []string{"/x", "=k=v"}))

	got := srv.AssertReceived(t, MatchCommand("/x"), "/x command")
	assert.Equal(t, []string{"/x", "=k=v"}, got)

	all := srv.AssertReceivedAll(t, MatchCommand("/x"), 1)
	assert.Len(t, all, 1)

	srv.AssertNotReceived(t, MatchCommand("/y"))
}

func TestServer_fallback_FIFO(t *testing.T) {
	// Existing test pattern: tanpa handler, Script tetap berfungsi.
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	srv.Script([]string{"!done", "=fifo=1"})

	conn, r := dialClient(t, srv)
	require.NoError(t, EncodeSentence(conn, []string{"/some-unknown-cmd"}))

	got, err := DecodeSentence(r)
	require.NoError(t, err)
	assert.Equal(t, []string{"!done", "=fifo=1"}, got)
}

func TestServer_handlersPrecedeFIFO(t *testing.T) {
	srv, err := Start()
	require.NoError(t, err)
	defer srv.Close()

	srv.OnSentence(MatchCommand("/known"), DoneReply("=src=handler"))
	srv.Script([]string{"!done", "=src=fifo"})

	conn, r := dialClient(t, srv)
	require.NoError(t, EncodeSentence(conn, []string{"/known"}))

	got, err := DecodeSentence(r)
	require.NoError(t, err)
	assert.Equal(t, []string{"!done", "=src=handler"}, got)
}

func TestBuilders_shapes(t *testing.T) {
	assert.Equal(t, []string{"!done"}, DoneReply())
	assert.Equal(t, []string{"!done", "=ret=*1"}, DoneReply("=ret=*1"))
	assert.Equal(t, []string{"!re", "=k=v"}, ReReply("=k=v"))
	assert.Equal(t, []string{"!trap", "=message=oops"}, TrapReply("oops", ""))
	assert.Equal(t, []string{"!trap", "=message=oops", "=category=0"}, TrapReply("oops", "0"))
	assert.Equal(t, []string{
		"!re", "=.id=*1", "=user=alice", "=address=1.2.3.4", "=mac-address=aa:bb:cc:dd:ee:ff",
	}, ActiveLoggedIn("1", "alice", "1.2.3.4", "aa:bb:cc:dd:ee:ff"))
	assert.Equal(t, []string{"!re", "=.id=*1", "=.dead=true"}, ActiveLoggedOut("*1"))

	want := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	row := UserPrintExpired("1", "user1", "default", want)
	assert.Equal(t, "=comment=jan/02/2026 15:04:05", row[4])
}
