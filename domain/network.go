package domain

// QueueSimple = row /queue/simple/print.
type QueueSimple struct {
	ID          string
	Name        string
	Target      string
	MaxLimit    string // mis. "1M/2M"
	LimitAt     string
	BurstLimit  string
	Parent      string
	Priority    string // "8/8" format
	BucketSize  string
	PacketMarks string
	Queue       string // queue type
	Disabled    bool
	Dynamic     bool
	Comment     string

	// Stats fields (di-isi saat print biasa, format string "in/out").
	Bytes              string
	Packets            string
	Rate               string
	TotalRate          string
	PacketRate         string
	TotalPacketRate    string
	QueuedBytes        string
	TotalQueuedBytes   string
	QueuedPackets      string
	TotalQueuedPackets string
	TotalBytes         string
	TotalPackets       string
	Dropped            string
	TotalDropped       string
}

// QueueSimpleWithStats = row /queue/simple/print stats interval=<d> (analisis §1.10).
type QueueSimpleWithStats = QueueSimple

// IPPool = row /ip/pool/print.
type IPPool struct {
	ID        string
	Name      string
	Ranges    string
	Total     int64
	Used      int64
	Available int64
	NextPool  string
	Comment   string
}

// ARPEntry = row /ip/arp/print.
type ARPEntry struct {
	ID         string
	Address    string
	MACAddress string
	Interface  string
	Dynamic    bool
	Disabled   bool
	Complete   bool
	Published  bool
	Invalid    bool
	Comment    string
}

// DHCPLease = row /ip/dhcp-server/lease/print.
type DHCPLease struct {
	ID           string
	Address      string
	MACAddress   string
	ClientID     string
	HostName     string
	Server       string
	Status       string
	ExpiresAfter string
	LastSeen     string
	Dynamic      bool
	Disabled     bool
	Comment      string
}

// Interface = row /interface/print.
type Interface struct {
	ID               string
	Name             string
	DefaultName      string
	Type             string
	MTU              string // bisa "auto" atau angka, simpan sebagai string
	ActualMTU        int64
	MACAddress       string
	LastLinkUpTime   string
	LastLinkDownTime string
	LinkDowns        int64
	Running          bool
	Disabled         bool
	Comment          string
}

// TrafficSnapshot = hasil /interface/monitor-traffic once="".
type TrafficSnapshot struct {
	Name           string
	RxBitsPerSec   int64
	TxBitsPerSec   int64
	RxPacketPerSec int64
	TxPacketPerSec int64
}

// PPPSecret = row /ppp/secret/print.
//
// Cross-ref: analisis §1.12 (inferred).
type PPPSecret struct {
	ID               string
	Name             string
	Password         string
	Service          string
	Profile          string
	LocalAddr        string
	RemoteAddr       string
	CallerID         string // caller-id: filter by calling station (mis. MAC untuk PPPoE)
	Routes           string
	IPv6Routes       string
	RemoteIPv6Prefix string
	LimitBytesIn     int64
	LimitBytesOut    int64
	LastLoggedOut        string
	LastCallerID         string
	LastDisconnectReason string
	Disabled             bool
	Comment              string
}

// PPPProfile = row /ppp/profile/print.
type PPPProfile struct {
	ID             string
	Name           string
	LocalAddr      string
	RemoteAddr     string
	RateLimit      string
	DNSServer      string
	Bridge         string
	ParentQueue    string
	IdleTimeout    string
	SessionTimeout string
	OnUp           string
	OnDown         string
	OnlyOne        string // "default" | "yes" | "no"
	UseCompression string // "default" | "yes" | "no"
	UseEncryption  string // "default" | "yes" | "no"
	ChangeTCPMSS   string // "default" | "yes" | "no"
	Disabled       bool
	Comment        string
}

// PPPActive = row /ppp/active/print.
type PPPActive struct {
	ID            string
	Name          string
	Service       string
	CallerID      string
	Address       string
	Uptime        string
	Encoding      string
	SessionID     string
	LimitBytesIn  int64
	LimitBytesOut int64
	Comment       string
}
