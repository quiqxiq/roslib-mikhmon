package hotspot

import (
	"context"
	"time"

	"github.com/quiqxiq/roslib"
)

const serverPath = "/ip/hotspot"

// HotspotServer adalah row /ip/hotspot/print (daftar nama server hotspot).
type HotspotServer struct {
	ID                string
	Name              string
	Profile           string
	Interface         string
	AddressPool       string
	AddressesPerMAC   string // RouterOS field "addresses-per-mac" (numeric atau "unlimited")
	IdleTimeout       string
	KeepaliveTimeout  string
	LoginTimeout      string
	Disabled          bool
}

// ServerList → /ip/hotspot/print (analisis §1.8).
func (c *Client) ServerList(ctx context.Context) ([]HotspotServer, error) {
	reply, err := c.dev.Path(serverPath).Print().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return sentencesToServers(reply.Rows), nil
}

// ServerListCached → /ip/hotspot/print dengan TTL cache. Server jarang
// berubah; cocok untuk dashboard dropdown yang refresh sering.
func (c *Client) ServerListCached(ctx context.Context, ttl time.Duration) ([]HotspotServer, error) {
	reply, err := c.dev.Path(serverPath).Print().ExecCached(ctx, ttl)
	if err != nil {
		return nil, err
	}
	return sentencesToServers(reply.Rows), nil
}

func sentencesToServers(rows []*roslib.Sentence) []HotspotServer {
	out := make([]HotspotServer, 0, len(rows))
	for _, s := range rows {
		out = append(out, HotspotServer{
			ID:               s.Get(".id"),
			Name:             s.Get("name"),
			Profile:          s.Get("profile"),
			Interface:        s.Get("interface"),
			AddressPool:      s.Get("address-pool"),
			AddressesPerMAC:  s.Get("addresses-per-mac"),
			IdleTimeout:      s.Get("idle-timeout"),
			KeepaliveTimeout: s.Get("keepalive-timeout"),
			LoginTimeout:     s.Get("login-timeout"),
			Disabled:         s.BoolOr("disabled", false),
		})
	}
	return out
}
