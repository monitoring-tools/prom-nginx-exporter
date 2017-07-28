package scraper_test

import (
	"strings"
	"testing"

	"github.com/monitoring-tools/prom-nginx-exporter/metric"
	"github.com/monitoring-tools/prom-nginx-exporter/scraper"
	. "gopkg.in/check.v1"
)

func TestNginxPlusScraper(t *testing.T) { TestingT(t) }

type NginxPlusScraperSuite struct{}

var _ = Suite(&NginxPlusScraperSuite{})

var validNginxPlusStats = `
{
    "version": 6,
    "nginx_version":  "1.22.333",
    "address":        "1.2.3.4",
    "generation":     88,
    "load_timestamp": 1451606400000,
    "timestamp":      1451606400000,
    "pid":            9999,
    "processes": {
        "respawned": 9999
     },
    "connections": {
        "accepted": 1234567890000,
        "dropped":  2345678900000,
        "active":   345,
        "idle":     567
    },
    "ssl": {
        "handshakes":        1234567800000,
        "handshakes_failed": 5432100000000,
        "session_reuses":    6543210000000
    },
    "requests": {
        "total":   9876543210000,
        "current": 98
    },
    "server_zones": {
        "zone.a_80": {
            "processing": 12,
            "requests": 34,
            "responses": {
                "1xx": 111,
                "2xx": 222,
                "3xx": 333,
                "4xx": 444,
                "5xx": 555,
                "total": 999
            },
            "discarded": 11,
            "received": 22,
            "sent": 33
        }
    },
    "upstreams": {
        "first_upstream": {
        	"queue": {
        		"size": 100,
        		"max_size": 1000,
        		"overflows": 12
			},
            "peers": [
                {
                    "id": 0,
                    "server": "1.2.3.123:80",
                    "backup": false,
                    "weight": 1,
                    "state": "up",
                    "active": 0,
                    "requests": 9876,
                    "responses": {
                        "1xx": 1111,
                        "2xx": 2222,
                        "3xx": 3333,
                        "4xx": 4444,
                        "5xx": 5555,
                        "total": 987654
                    },
                    "sent": 987654321,
                    "received": 87654321,
                    "fails": 98,
                    "unavail": 65,
                    "health_checks": {
                        "checks": 54,
                        "fails": 32,
                        "unhealthy": 21,
                        "last_passed": false
                    },
                    "downtime": 5432,
                    "downstart": 4321,
                    "selected": 1451606400000,
                    "header_time": 2451606400000,
                    "response_time": 3451606400000,
                    "max_conns": 1000000
                }
            ],
            "keepalive": 1,
            "zombies": 2
        }
    },
    "caches": {
        "cache_01": {
            "size": 12,
            "max_size": 23,
            "cold": false,
            "hit": {
                "responses": 34,
                "bytes": 45
            },
            "stale": {
                "responses": 56,
                "bytes": 67
            },
            "updating": {
                "responses": 78,
                "bytes": 89
            },
            "revalidated": {
                "responses": 90,
                "bytes": 98
            },
            "miss": {
                "responses": 87,
                "bytes": 76,
                "responses_written": 65,
                "bytes_written": 54
            },
            "expired": {
                "responses": 43,
                "bytes": 32,
                "responses_written": 21,
                "bytes_written": 10
            },
            "bypass": {
                "responses": 13,
                "bytes": 35,
                "responses_written": 57,
                "bytes_written": 79
            }
        }
    },
    "stream": {
        "server_zones": {
            "stream.zone.01": {
                "processing": 24,
                "connections": 46,
                "received": 68,
                "sent": 80
            }
        },
        "upstreams": {
            "upstream.01": {
                "peers": [
                    {
                        "id": 1,
                        "server": "5.4.3.2:2345",
                        "backup": false,
                        "weight": 1,
                        "state": "up",
                        "active": 0,
                        "connections": 0,
                        "sent": 0,
                        "received": 0,
                        "fails": 0,
                        "unavail": 0,
                        "downtime": 0,
                        "downstart": 0,
                        "selected": 0,
                        "health_checks": {
                            "checks": 40851,
                            "fails": 0,
                            "unhealthy": 0,
                            "last_passed": true
                        },
                        "connect_time": 993,
                        "first_byte_time": 994,
                        "response_time": 995
                    }
                ],
                "zombies": 0
            }
        }
    }
}
`

func (s NginxPlusScraperSuite) TestScrape_Success(c *C) {
	nginxPlusScraper := scraper.NewNginxPlusScraper()
	reader := strings.NewReader(validNginxPlusStats)

	metrics := make(chan metric.Metric, 98)
	labels := map[string]string{
		"host": "zone.a_80",
		"port": "8080",
	}

	err := nginxPlusScraper.Scrape(reader, metrics, labels)
	c.Assert(err, IsNil, Commentf("error occurred during scrape nginx plus stats"))

	m := <-metrics
	c.Assert(m.Name, Equals, "processes_respawned", Commentf("incorrect metrics name of 'processes_respawned' field"))
	c.Assert(m.Value, Equals, 9999, Commentf("incorrect value of metric 'processes_respawned'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "connections_accepted", Commentf("incorrect metrics name of 'connections_accepted' field"))
	c.Assert(m.Value, Equals, 1234567890000, Commentf("incorrect value of metric 'connections_accepted'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "connections_dropped", Commentf("incorrect metrics name of 'connections_dropped' field"))
	c.Assert(m.Value, Equals, 2345678900000, Commentf("incorrect value of metric 'connections_dropped'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "connections_active", Commentf("incorrect metrics name of 'connections_active' field"))
	c.Assert(m.Value, Equals, 345, Commentf("incorrect value of metric 'connections_active'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "connections_idle", Commentf("incorrect metrics name of 'connections_idle' field"))
	c.Assert(m.Value, Equals, 567, Commentf("incorrect value of metric 'connections_idle'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "ssl_handshakes", Commentf("incorrect metrics name of 'ssl_handshakes' field"))
	c.Assert(m.Value, Equals, int64(1234567800000), Commentf("incorrect value of metric 'ssl_handshakes'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "ssl_handshakes_failed", Commentf("incorrect metrics name of 'ssl_handshakes_failed' field"))
	c.Assert(m.Value, Equals, int64(5432100000000), Commentf("incorrect value of metric 'ssl_handshakes_failed'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "ssl_session_reuses", Commentf("incorrect metrics name of 'ssl_session_reuses' field"))
	c.Assert(m.Value, Equals, int64(6543210000000), Commentf("incorrect value of metric 'ssl_session_reuses'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "requests_total", Commentf("incorrect metrics name of 'requests_total' field"))
	c.Assert(m.Value, Equals, int64(9876543210000), Commentf("incorrect value of metric 'requests_total'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "requests_current", Commentf("incorrect metrics name of 'requests_current' field"))
	c.Assert(m.Value, Equals, int(98), Commentf("incorrect value of metric 'requests_current'"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	zoneLabels := make(map[string]string)
	for l, v := range labels {
		zoneLabels[l] = v
	}
	zoneLabels["zone"] = "zone.a_80"

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_processing", Commentf("incorrect metrics name of 'zone_processing' field"))
	c.Assert(m.Value, Equals, 12, Commentf("incorrect value of metric 'zone_processing'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_requests", Commentf("incorrect metrics name of 'zone_requests' field"))
	c.Assert(m.Value, Equals, int64(34), Commentf("incorrect value of metric 'zone_requests'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_responses_1xx", Commentf("incorrect metrics name of 'zone_responses_1xx' field"))
	c.Assert(m.Value, Equals, int64(111), Commentf("incorrect value of metric 'zone_responses_1xx'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_responses_2xx", Commentf("incorrect metrics name of 'zone_responses_2xx' field"))
	c.Assert(m.Value, Equals, int64(222), Commentf("incorrect value of metric 'zone_responses_2xx'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_responses_3xx", Commentf("incorrect metrics name of 'zone_responses_3xx' field"))
	c.Assert(m.Value, Equals, int64(333), Commentf("incorrect value of metric 'zone_responses_3xx'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_responses_4xx", Commentf("incorrect metrics name of 'zone_responses_4xx' field"))
	c.Assert(m.Value, Equals, int64(444), Commentf("incorrect value of metric 'zone_responses_4xx'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_responses_5xx", Commentf("incorrect metrics name of 'zone_responses_5xx' field"))
	c.Assert(m.Value, Equals, int64(555), Commentf("incorrect value of metric 'zone_responses_5xx'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_responses_total", Commentf("incorrect metrics name of 'zone_responses_total' field"))
	c.Assert(m.Value, Equals, int64(999), Commentf("incorrect value of metric 'zone_responses_total'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_received", Commentf("incorrect metrics name of 'zone_received' field"))
	c.Assert(m.Value, Equals, int64(22), Commentf("incorrect value of metric 'zone_received'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_sent", Commentf("incorrect metrics name of 'zone_sent' field"))
	c.Assert(m.Value, Equals, int64(33), Commentf("incorrect value of metric 'zone_sent'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "zone_discarded", Commentf("incorrect metrics name of 'zone_discarded' field"))
	c.Assert(m.Value, Equals, int64(11), Commentf("incorrect value of metric 'zone_discarded'"))
	c.Assert(m.Labels, DeepEquals, zoneLabels, Commentf("incorrect set of labels"))

	upstramLabels := make(map[string]string)
	for l, v := range labels {
		upstramLabels[l] = v
	}
	upstramLabels["upstream"] = "first_upstream"

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_keepalive", Commentf("incorrect metrics name of 'upstream_keepalive' field"))
	c.Assert(m.Value, Equals, 1, Commentf("incorrect value of metric 'upstream_keepalive'"))
	c.Assert(m.Labels, DeepEquals, upstramLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_zombies", Commentf("incorrect metrics name of 'upstream_zombies' field"))
	c.Assert(m.Value, Equals, 2, Commentf("incorrect value of metric 'upstream_zombies'"))
	c.Assert(m.Labels, DeepEquals, upstramLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_queue_size", Commentf("incorrect metrics name of 'upstream_queue_size' field"))
	c.Assert(m.Value, Equals, 100, Commentf("incorrect value of metric 'upstream_queue_size'"))
	c.Assert(m.Labels, DeepEquals, upstramLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_queue_max_size", Commentf("incorrect metrics name of 'upstream_queue_max_size' field"))
	c.Assert(m.Value, Equals, 1000, Commentf("incorrect value of metric 'upstream_queue_max_size'"))
	c.Assert(m.Labels, DeepEquals, upstramLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_queue_overflows", Commentf("incorrect metrics name of 'upstream_queue_overflows' field"))
	c.Assert(m.Value, Equals, int64(12), Commentf("incorrect value of metric 'upstream_queue_overflows'"))
	c.Assert(m.Labels, DeepEquals, upstramLabels, Commentf("incorrect set of labels"))

	peerLabels := make(map[string]string)
	for l, v := range upstramLabels {
		peerLabels[l] = v
	}
	peerLabels["id"] = "0"
	peerLabels["serverAddress"] = "1.2.3.123:80"

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_backup", Commentf("incorrect metrics name of 'upstream_peer_backup' field"))
	c.Assert(m.Value, Equals, false, Commentf("incorrect value of metric 'upstream_peer_backup'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_weight", Commentf("incorrect metrics name of 'upstream_peer_weight' field"))
	c.Assert(m.Value, Equals, 1, Commentf("incorrect value of metric 'upstream_peer_weight'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_state", Commentf("incorrect metrics name of 'upstream_peer_state' field"))
	c.Assert(m.Value, Equals, "up", Commentf("incorrect value of metric 'upstream_peer_state'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_active", Commentf("incorrect metrics name of 'upstream_peer_active' field"))
	c.Assert(m.Value, Equals, 0, Commentf("incorrect value of metric 'upstream_peer_active'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_requests", Commentf("incorrect metrics name of 'upstream_peer_requests' field"))
	c.Assert(m.Value, Equals, int64(9876), Commentf("incorrect value of metric 'upstream_peer_requests'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_responses_1xx", Commentf("incorrect metrics name of 'upstream_peer_responses_1xx' field"))
	c.Assert(m.Value, Equals, int64(1111), Commentf("incorrect value of metric 'upstream_peer_responses_1xx'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_responses_2xx", Commentf("incorrect metrics name of 'upstream_peer_responses_2xx' field"))
	c.Assert(m.Value, Equals, int64(2222), Commentf("incorrect value of metric 'upstream_peer_responses_2xx'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_responses_3xx", Commentf("incorrect metrics name of 'upstream_peer_responses_3xx' field"))
	c.Assert(m.Value, Equals, int64(3333), Commentf("incorrect value of metric 'upstream_peer_responses_3xx'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_responses_4xx", Commentf("incorrect metrics name of 'upstream_peer_responses_4xx' field"))
	c.Assert(m.Value, Equals, int64(4444), Commentf("incorrect value of metric 'upstream_peer_responses_4xx'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_responses_5xx", Commentf("incorrect metrics name of 'upstream_peer_responses_5xx' field"))
	c.Assert(m.Value, Equals, int64(5555), Commentf("incorrect value of metric 'upstream_peer_responses_5xx'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_responses_total", Commentf("incorrect metrics name of 'upstream_peer_responses_total' field"))
	c.Assert(m.Value, Equals, int64(987654), Commentf("incorrect value of metric 'upstream_peer_responses_total'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_sent", Commentf("incorrect metrics name of 'upstream_peer_sent' field"))
	c.Assert(m.Value, Equals, int64(987654321), Commentf("incorrect value of metric 'upstream_peer_sent'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_received", Commentf("incorrect metrics name of 'upstream_peer_received' field"))
	c.Assert(m.Value, Equals, int64(87654321), Commentf("incorrect value of metric 'upstream_peer_received'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_fails", Commentf("incorrect metrics name of 'upstream_peer_fails' field"))
	c.Assert(m.Value, Equals, int64(98), Commentf("incorrect value of metric 'upstream_peer_fails'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_unavail", Commentf("incorrect metrics name of 'upstream_peer_unavail' field"))
	c.Assert(m.Value, Equals, int64(65), Commentf("incorrect value of metric 'upstream_peer_unavail'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_healthchecks_checks", Commentf("incorrect metrics name of 'upstream_peer_healthchecks_checks' field"))
	c.Assert(m.Value, Equals, int64(54), Commentf("incorrect value of metric 'upstream_peer_healthchecks_checks'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_healthchecks_fails", Commentf("incorrect metrics name of 'upstream_peer_healthchecks_fails' field"))
	c.Assert(m.Value, Equals, int64(32), Commentf("incorrect value of metric 'upstream_peer_healthchecks_fails'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_healthchecks_unhealthy", Commentf("incorrect metrics name of 'upstream_peer_healthchecks_unhealthy' field"))
	c.Assert(m.Value, Equals, int64(21), Commentf("incorrect value of metric 'upstream_peer_healthchecks_unhealthy'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_downtime", Commentf("incorrect metrics name of 'upstream_peer_downtime' field"))
	c.Assert(m.Value, Equals, int64(5432), Commentf("incorrect value of metric 'upstream_peer_downtime'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_downstart", Commentf("incorrect metrics name of 'upstream_peer_downstart' field"))
	c.Assert(m.Value, Equals, int64(4321), Commentf("incorrect value of metric 'upstream_peer_downstart'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_selected", Commentf("incorrect metrics name of 'upstream_peer_selected' field"))
	c.Assert(m.Value, Equals, int64(1451606400000), Commentf("incorrect value of metric 'upstream_peer_selected'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_healthchecks_last_passed", Commentf("incorrect metrics name of 'upstream_peer_healthchecks_last_passed' field"))
	c.Assert(m.Value, Equals, false, Commentf("incorrect value of metric 'upstream_peer_healthchecks_last_passed'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_header_time", Commentf("incorrect metrics name of 'upstream_peer_header_time' field"))
	c.Assert(m.Value, Equals, int64(2451606400000), Commentf("incorrect value of metric 'upstream_peer_header_time'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_response_time", Commentf("incorrect metrics name of 'upstream_peer_response_time' field"))
	c.Assert(m.Value, Equals, int64(3451606400000), Commentf("incorrect value of metric 'upstream_peer_response_time'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "upstream_peer_max_conns", Commentf("incorrect metrics name of 'upstream_peer_max_conns' field"))
	c.Assert(m.Value, Equals, 1000000, Commentf("incorrect value of metric 'upstream_peer_max_conns'"))
	c.Assert(m.Labels, DeepEquals, peerLabels, Commentf("incorrect set of labels"))

	cacheLabels := make(map[string]string)
	for l, v := range labels {
		cacheLabels[l] = v
	}
	cacheLabels["cache"] = "cache_01"

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_size", Commentf("incorrect metrics name of 'cache_size' field"))
	c.Assert(m.Value, Equals, int64(12), Commentf("incorrect value of metric 'cache_size'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_max_size", Commentf("incorrect metrics name of 'cache_max_size' field"))
	c.Assert(m.Value, Equals, int64(12), Commentf("incorrect value of metric 'cache_max_size'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_cold", Commentf("incorrect metrics name of 'cache_cold' field"))
	c.Assert(m.Value, Equals, false, Commentf("incorrect value of metric 'cache_cold'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_hit_responses", Commentf("incorrect metrics name of 'cache_hit_responses' field"))
	c.Assert(m.Value, Equals, int64(34), Commentf("incorrect value of metric 'cache_hit_responses'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_hit_bytes", Commentf("incorrect metrics name of 'cache_hit_bytes' field"))
	c.Assert(m.Value, Equals, int64(45), Commentf("incorrect value of metric 'cache_hit_bytes'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_stale_responses", Commentf("incorrect metrics name of 'cache_stale_responses' field"))
	c.Assert(m.Value, Equals, int64(56), Commentf("incorrect value of metric 'cache_stale_responses'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_stale_bytes", Commentf("incorrect metrics name of 'cache_stale_bytes' field"))
	c.Assert(m.Value, Equals, int64(67), Commentf("incorrect value of metric 'cache_stale_bytes'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_updating_responses", Commentf("incorrect metrics name of 'cache_updating_responses' field"))
	c.Assert(m.Value, Equals, int64(78), Commentf("incorrect value of metric 'cache_updating_responses'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_updating_bytes", Commentf("incorrect metrics name of 'cache_updating_bytes' field"))
	c.Assert(m.Value, Equals, int64(89), Commentf("incorrect value of metric 'cache_updating_bytes'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_revalidated_responses", Commentf("incorrect metrics name of 'cache_revalidated_responses' field"))
	c.Assert(m.Value, Equals, int64(90), Commentf("incorrect value of metric 'cache_revalidated_responses'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_revalidated_bytes", Commentf("incorrect metrics name of 'cache_revalidated_bytes' field"))
	c.Assert(m.Value, Equals, int64(98), Commentf("incorrect value of metric 'cache_revalidated_bytes'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_miss_responses", Commentf("incorrect metrics name of 'cache_miss_responses' field"))
	c.Assert(m.Value, Equals, int64(87), Commentf("incorrect value of metric 'cache_miss_responses'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_miss_bytes", Commentf("incorrect metrics name of 'cache_miss_bytes' field"))
	c.Assert(m.Value, Equals, int64(76), Commentf("incorrect value of metric 'cache_miss_bytes'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_miss_responses_written", Commentf("incorrect metrics name of 'cache_miss_responses_written' field"))
	c.Assert(m.Value, Equals, int64(65), Commentf("incorrect value of metric 'cache_miss_responses_written'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_miss_bytes_written", Commentf("incorrect metrics name of 'cache_miss_bytes_written' field"))
	c.Assert(m.Value, Equals, int64(76), Commentf("incorrect value of metric 'cache_miss_bytes_written'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_expired_responses", Commentf("incorrect metrics name of 'cache_expired_responses' field"))
	c.Assert(m.Value, Equals, int64(43), Commentf("incorrect value of metric 'cache_expired_responses'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_expired_bytes", Commentf("incorrect metrics name of 'cache_expired_bytes' field"))
	c.Assert(m.Value, Equals, int64(32), Commentf("incorrect value of metric 'cache_expired_bytes'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_expired_responses_written", Commentf("incorrect metrics name of 'cache_expired_responses_written' field"))
	c.Assert(m.Value, Equals, int64(21), Commentf("incorrect value of metric 'cache_expired_responses_written'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_expired_bytes_written", Commentf("incorrect metrics name of 'cache_expired_bytes_written' field"))
	c.Assert(m.Value, Equals, int64(10), Commentf("incorrect value of metric 'cache_expired_bytes_written'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_responses", Commentf("incorrect metrics name of 'cache_responses' field"))
	c.Assert(m.Value, Equals, int64(13), Commentf("incorrect value of metric 'cache_responses'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_bytes", Commentf("incorrect metrics name of 'cache_bytes' field"))
	c.Assert(m.Value, Equals, int64(35), Commentf("incorrect value of metric 'cache_bytes'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_responses_written", Commentf("incorrect metrics name of 'cache_responses_written' field"))
	c.Assert(m.Value, Equals, int64(57), Commentf("incorrect value of metric 'cache_responses_written'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "cache_bytes_written", Commentf("incorrect metrics name of 'cache_bytes_written' field"))
	c.Assert(m.Value, Equals, int64(79), Commentf("incorrect value of metric 'cache_bytes_written'"))
	c.Assert(m.Labels, DeepEquals, cacheLabels, Commentf("incorrect set of labels"))

	streamLabels := make(map[string]string)
	for l, v := range labels {
		streamLabels[l] = v
	}
	streamLabels["zone"] = "stream.zone.01"

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_zone_processing", Commentf("incorrect metrics name of 'stream_zone_processing' field"))
	c.Assert(m.Value, Equals, 24, Commentf("incorrect value of metric 'stream_zone_processing'"))
	c.Assert(m.Labels, DeepEquals, streamLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_zone_connections", Commentf("incorrect metrics name of 'stream_zone_connections' field"))
	c.Assert(m.Value, Equals, 46, Commentf("incorrect value of metric 'stream_zone_connections'"))
	c.Assert(m.Labels, DeepEquals, streamLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_zone_received", Commentf("incorrect metrics name of 'stream_zone_received' field"))
	c.Assert(m.Value, Equals, int64(68), Commentf("incorrect value of metric 'stream_zone_received'"))
	c.Assert(m.Labels, DeepEquals, streamLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_zone_sent", Commentf("incorrect metrics name of 'stream_zone_sent' field"))
	c.Assert(m.Value, Equals, int64(80), Commentf("incorrect value of metric 'stream_zone_sent'"))
	c.Assert(m.Labels, DeepEquals, streamLabels, Commentf("incorrect set of labels"))

	streamUpstreamLabels := make(map[string]string)
	for l, v := range labels {
		streamUpstreamLabels[l] = v
	}
	streamUpstreamLabels["upstream"] = "upstream.01"

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_zombies", Commentf("incorrect metrics name of 'stream_upstream_zombies' field"))
	c.Assert(m.Value, Equals, 0, Commentf("incorrect value of metric 'stream_upstream_zombies'"))
	c.Assert(m.Labels, DeepEquals, streamUpstreamLabels, Commentf("incorrect set of labels"))

	streamPeerLabels := make(map[string]string)
	for l, v := range streamUpstreamLabels {
		streamPeerLabels[l] = v
	}
	streamPeerLabels["serverAddress"] = "5.4.3.2:2345"
	streamPeerLabels["id"] = "1"

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_backup", Commentf("incorrect metrics name of 'stream_upstream_peer_backup' field"))
	c.Assert(m.Value, Equals, false, Commentf("incorrect value of metric 'stream_upstream_peer_backup'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_weight", Commentf("incorrect metrics name of 'stream_upstream_peer_weight' field"))
	c.Assert(m.Value, Equals, 1, Commentf("incorrect value of metric 'stream_upstream_peer_weight'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_state", Commentf("incorrect metrics name of 'stream_upstream_peer_state' field"))
	c.Assert(m.Value, Equals, "up", Commentf("incorrect value of metric 'stream_upstream_peer_state'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_active", Commentf("incorrect metrics name of 'stream_upstream_peer_active' field"))
	c.Assert(m.Value, Equals, 0, Commentf("incorrect value of metric 'stream_upstream_peer_active'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_connections", Commentf("incorrect metrics name of 'stream_upstream_peer_connections' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_connections'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_sent", Commentf("incorrect metrics name of 'stream_upstream_peer_sent' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_sent'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_received", Commentf("incorrect metrics name of 'stream_upstream_peer_received' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_received'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_fails", Commentf("incorrect metrics name of 'stream_upstream_peer_fails' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_fails'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_unavail", Commentf("incorrect metrics name of 'stream_upstream_peer_unavail' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_unavail'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_healthchecks_checks", Commentf("incorrect metrics name of 'stream_upstream_peer_healthchecks_checks' field"))
	c.Assert(m.Value, Equals, int64(40851), Commentf("incorrect value of metric 'stream_upstream_peer_healthchecks_checks'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_healthchecks_fails", Commentf("incorrect metrics name of 'stream_upstream_peer_healthchecks_fails' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_healthchecks_fails'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_healthchecks_unhealthy", Commentf("incorrect metrics name of 'stream_upstream_peer_healthchecks_unhealthy' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_healthchecks_unhealthy'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_healthchecks_downtime", Commentf("incorrect metrics name of 'stream_upstream_peer_healthchecks_downtime' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_healthchecks_downtime'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_healthchecks_downstart", Commentf("incorrect metrics name of 'stream_upstream_peer_healthchecks_downstart' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_healthchecks_downstart'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_healthchecks_selected", Commentf("incorrect metrics name of 'stream_upstream_peer_healthchecks_selected' field"))
	c.Assert(m.Value, Equals, int64(0), Commentf("incorrect value of metric 'stream_upstream_peer_healthchecks_selected'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_healthchecks_last_passed", Commentf("incorrect metrics name of 'stream_upstream_peer_healthchecks_last_passed' field"))
	c.Assert(m.Value, Equals, true, Commentf("incorrect value of metric 'stream_upstream_peer_healthchecks_last_passed'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_connect_time", Commentf("incorrect metrics name of 'stream_upstream_peer_connect_time' field"))
	c.Assert(m.Value, Equals, 993, Commentf("incorrect value of metric 'stream_upstream_peer_connect_time'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_first_byte_time", Commentf("incorrect metrics name of 'stream_upstream_peer_first_byte_time' field"))
	c.Assert(m.Value, Equals, 994, Commentf("incorrect value of metric 'stream_upstream_peer_first_byte_time'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))

	m = <-metrics
	c.Assert(m.Name, Equals, "stream_upstream_peer_response_time", Commentf("incorrect metrics name of 'stream_upstream_peer_response_time' field"))
	c.Assert(m.Value, Equals, 995, Commentf("incorrect value of metric 'stream_upstream_peer_response_time'"))
	c.Assert(m.Labels, DeepEquals, streamPeerLabels, Commentf("incorrect set of labels"))
}

func (s NginxPlusScraperSuite) TestScrape_Fail(c *C) {
	nginxPlusScraper := scraper.NewNginxPlusScraper()
	reader := strings.NewReader(`{"version":"invalid json"}`)

	metrics := make(chan metric.Metric, 96)
	labels := map[string]string{"host": "zone.a_80", "port": "8080"}

	err := nginxPlusScraper.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("error should be occurred"))
	c.Assert(err.Error(), Equals, "Error while decoding JSON response", Commentf("incorrect error massage of parsing json"))
}
