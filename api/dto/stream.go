package dto

// InterfaceStatsEvent adalah typed event untuk SSE /stream/network/interfaces/stats.
type InterfaceStatsEvent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	RxByte   int64  `json:"rx_byte"`
	TxByte   int64  `json:"tx_byte"`
	RxPacket int64  `json:"rx_packet"`
	TxPacket int64  `json:"tx_packet"`
	Running  bool   `json:"running"`
	Disabled bool   `json:"disabled"`
}

// QueueStatsEvent adalah typed event untuk SSE /stream/network/queues/stats.
type QueueStatsEvent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Target   string `json:"target"`
	Bytes    string `json:"bytes"`   // format RouterOS "in/out"
	Packets  string `json:"packets"` // format RouterOS "in/out"
	Rate     string `json:"rate"`
	MaxLimit string `json:"max_limit,omitempty"`
}

// HotspotActiveEvent adalah typed event untuk SSE /stream/hotspot/active.
// Field dead=true menandakan session yang terminate (event dari print follow).
type HotspotActiveEvent struct {
	ID         string `json:"id"`
	User       string `json:"user"`
	Address    string `json:"address,omitempty"`
	MACAddress string `json:"mac_address,omitempty"`
	Server     string `json:"server,omitempty"`
	LoginBy    string `json:"login_by,omitempty"`
	BytesIn    int64  `json:"bytes_in"`
	BytesOut   int64  `json:"bytes_out"`
	Uptime     string `json:"uptime,omitempty"`
	Dead       bool   `json:"dead"`
}

// PPPActiveEvent adalah typed event untuk SSE /stream/ppp/active.
// Field dead=true menandakan session yang terminate (event dari print follow).
type PPPActiveEvent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Service  string `json:"service,omitempty"`
	CallerID string `json:"caller_id,omitempty"`
	Address  string `json:"address,omitempty"`
	Uptime   string `json:"uptime,omitempty"`
	Dead     bool   `json:"dead"`
}

// LogEvent adalah typed event untuk SSE /stream/log.
type LogEvent struct {
	ID      string `json:"id"`
	Time    string `json:"time"`
	Topics  string `json:"topics"`
	Message string `json:"message"`
}
