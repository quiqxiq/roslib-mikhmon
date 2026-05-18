package domain

import "time"

// SystemIdentity = /system/identity/print → .name.
type SystemIdentity struct {
	Name string
}

// SystemResource = /system/resource/print.
type SystemResource struct {
	Uptime               string
	Version              string
	BuildTime            string
	FactorySoftware      string
	BoardName            string
	Platform             string
	CPU                  string
	CPUCount             int
	CPUFrequency         int64
	CPULoad              int
	FreeMemory           int64
	TotalMemory          int64
	FreeHDDSpace         int64
	TotalHDDSpace        int64
	ArchitectureName     string
	WriteSectSinceReboot int64
	BadBlocks            string
}

// SystemRouterboard = /system/routerboard/print.
type SystemRouterboard struct {
	Routerboard     bool
	Model           string
	BoardName       string
	SerialNumber    string
	FirmwareType    string
	FactoryFirmware string
	CurrentFirmware string
	UpgradeFirmware string
}

// SystemClock = /system/clock/print.
type SystemClock struct {
	Time                string
	Date                string
	TimeZoneName        string
	GMTOffset           string
	TimeZoneAutodetect  bool
	DSTActive           bool
}

// LogEntry = row /log/print.
type LogEntry struct {
	ID      string
	Time    string
	Topics  string
	Message string
	When    time.Time // best-effort, ParseTime saat list (boleh zero)
}

// Scheduler = row /system/scheduler/print.
type Scheduler struct {
	ID        string
	Name      string
	StartDate string
	StartTime string
	Interval  string
	OnEvent   string
	NextRun   string
	Policy    string
	Disabled  bool
	Comment   string
	RunCount  int
}

// Script = row /system/script/print.
//
// Mikhmon menyimpan 2 jenis: transaksi (comment="mikhmon") dan
// QuickPrint config (comment="QuickPrintMikhmon").
type Script struct {
	ID                     string
	Name                   string
	Owner                  string
	Source                 string
	Comment                string
	Policy                 string
	RunCount               int
	LastStarted            string
	DontRequirePermissions bool
}
