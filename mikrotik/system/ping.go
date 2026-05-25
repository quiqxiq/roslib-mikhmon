package system

import (
	"strconv"

	"github.com/quiqxiq/roslib"
	"github.com/quiqxiq/roslib/stream"
)

// PingStream → /ping address=<addr> (inherent stream, emit tiap reply).
// Caller wajib panggil StopPingStream(id) saat selesai.
func (c *Client) PingStream(id, address string, h func(*roslib.Sentence)) error {
	return c.dev.Path("/ping").
		With("address", address).
		Stream(id, h)
}

// PingCount → /ping address=<addr> count=<count>.
// Inherently streaming: router emit !re tiap reply lalu !done saat count tercapai.
// Callback onFinish dipanggil saat natural completion (!done) atau error.
func (c *Client) PingCount(id, address string, count int, h func(*roslib.Sentence), onFinish stream.FinishCallback) error {
	return c.dev.Path("/ping").
		With("address", address).
		With("count", strconv.Itoa(count)).
		OnFinish(onFinish).
		Stream(id, h)
}

// StopPingStream menghentikan stream ping (kirim /cancel ke router).
// Return true bila stream ada dan dihapus.
func (c *Client) StopPingStream(id string) bool {
	return c.dev.UnregisterStream(id)
}
