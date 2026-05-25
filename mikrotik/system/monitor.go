package system

import (
	"time"

	"github.com/quiqxiq/roslib"
)

// MonitorResource → poll /system/resource/print tiap interval (analisis §1.1).
// Cocok untuk dashboard CPU/RAM/uptime feed. Proplist memperkecil payload.
func (c *Client) MonitorResource(id string, interval time.Duration, h func(*roslib.Sentence)) error {
	return c.dev.RegisterPoll(roslib.PollConfig{
		ID:       id,
		Path:     "/system/resource",
		Args:     []string{"print", "proplist=uptime,version,build-time,board-name,cpu,cpu-count,cpu-frequency,cpu-load,free-memory,total-memory,free-hdd-space,total-hdd-space,architecture-name,write-sect-since-reboot,bad-blocks"},
		Interval: interval,
		Handler:  h,
	})
}

// MonitorRouterboard → poll /system/routerboard/print tiap interval (analisis §1.1).
// Cocok untuk dashboard board info & firmware.
func (c *Client) MonitorRouterboard(id string, interval time.Duration, h func(*roslib.Sentence)) error {
	return c.dev.RegisterPoll(roslib.PollConfig{
		ID:       id,
		Path:     "/system/routerboard",
		Args:     []string{"print", "proplist=routerboard,board-name,model,revision,serial-number,firmware-type,current-firmware,upgrade-firmware"},
		Interval: interval,
		Handler:  h,
	})
}

// MonitorScheduler → poll /system/scheduler/print ?name=<name> tiap interval
// (analisis §1.4). Dipakai untuk mendeteksi expiry user — mikhmon membaca
// field next-run dari handler.
func (c *Client) MonitorScheduler(id, name string, interval time.Duration, h func(*roslib.Sentence)) error {
	return c.dev.RegisterPoll(roslib.PollConfig{
		ID:       id,
		Path:     schedulerPath,
		Args:     []string{"print"},
		Where:    []roslib.WherePair{roslib.Where("name", name)},
		Interval: interval,
		Handler:  h,
	})
}

// StopMonitor menghentikan poll dengan ID tersebut.
// Return true bila poll ada dan dihapus.
func (c *Client) StopMonitor(id string) bool {
	return c.dev.UnregisterPoll(id)
}
