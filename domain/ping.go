package domain

// PingResult = satu reply ICMP echo dari RouterOS.
type PingResult struct {
	Seq    int
	Host   string
	Size   int
	TTL    int
	TimeMs float64 // parsed dari "24ms" atau "84ms435us"
	Status string  // kosong = success, kalau ada = timeout/unreachable
}

// PingSummary = summary akhir dari command /ping count=N.
type PingSummary struct {
	Target            string
	Sent              int
	Received          int
	PacketLossPercent float64
	MinRttMs          float64
	AvgRttMs          float64
	MaxRttMs          float64
	Results           []PingResult
}
