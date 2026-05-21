package sse

import "fmt"

// Topic constants — naming stabil untuk monitoring/log.
const (
	TopicHotspotActive           = "hotspot-active"
	TopicHotspotActiveFollowOnly = "hotspot-active-follow-only"
	TopicPPPActive               = "ppp-active"

	// Tabel follow streams (analisis §1.6, §1.12).
	TopicHotspotUser           = "hotspot-user"
	TopicHotspotUserFollowOnly = "hotspot-user-follow-only"
	TopicPPPSecret             = "ppp-secret"
	TopicPPPSecretFollowOnly   = "ppp-secret-follow-only"

	// Derived inactive streams (analisis §4 — workflows).
	TopicHotspotInactive = "hotspot-inactive"
	TopicPPPInactive     = "ppp-inactive"
)

// TopicLog membentuk topic untuk subscription /log/print follow dengan
// filter topics. Kosong = semua topics.
func TopicLog(topics string) string {
	if topics == "" {
		return "log:all"
	}
	return "log:" + topics
}

// TopicResource membentuk topic untuk poll /system/resource interval.
func TopicResource(interval string) string {
	return fmt.Sprintf("resource:%s", interval)
}

// TopicInterfaceTraffic untuk inherent stream monitor-traffic per iface.
func TopicInterfaceTraffic(iface string) string {
	return "traffic:" + iface
}

// TopicInterfaceStats untuk /interface/print stats interval=...
func TopicInterfaceStats(interval string) string {
	return "interface-stats:" + interval
}

// TopicQueueStats untuk /queue/simple/print stats interval=...
func TopicQueueStats(interval string) string {
	return "queue-stats:" + interval
}
