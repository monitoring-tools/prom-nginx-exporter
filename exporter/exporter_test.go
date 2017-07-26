package exporter_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/monitoring-tools/prom-nginx-exporter/exporter"
	"github.com/monitoring-tools/prom-nginx-exporter/scraper"

	"github.com/prometheus/client_golang/prometheus"
	. "gopkg.in/check.v1"
)

func TestNginxExporter(t *testing.T) { TestingT(t) }

type NginxExporterSuite struct{}

var _ = Suite(&NginxExporterSuite{})

var nginxStats = "Active connections: 2\n server accepts handled requests\n8522429 8522429 8641727\nReading: 0 Writing: 1 Waiting: 3"
var nginxPlusStats = `
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

func (s NginxExporterSuite) TestNginxStatsScrape_Success(c *C) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/plain")
	response := http.Response{
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       NewDummyBody(nginxStats),
	}

	client := &http.Client{Transport: NewDummyTransport(response)}
	exp := exporter.NewNginxPlusExporter(
		client,
		scraper.NewNginxScraper(),
		scraper.NewNginxPlusScraper(),
		"nginx_test",
		[]string{"http://localhost:9000"},
		[]string{},
	)

	metrics := make(chan prometheus.Metric)

	go func() {
		exp.Collect(metrics)
		close(metrics)
	}()

	checks := map[string]bool{
		"nginx_test_exporter_scrapes_total": false,
		"nginx_test_requests":               false,
		"nginx_test_reading":                false,
		"nginx_test_writing":                false,
		"nginx_test_waiting":                false,
		"nginx_test_active":                 false,
		"nginx_test_accepts":                false,
		"nginx_test_handled":                false,
	}

	for m := range metrics {
		switch m.(type) {
		case prometheus.Gauge:
			for metricName := range checks {
				if strings.Contains(m.Desc().String(), metricName) {
					checks[metricName] = true
				}
			}
		}
	}

	for metricName, exists := range checks {
		if !exists {
			c.Errorf("didn't find metric '%s'", metricName)
		}
	}
}

func (s NginxExporterSuite) TestNginxPlusStatsScrape_Success(c *C) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	response := http.Response{
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       NewDummyBody(nginxPlusStats),
	}

	client := &http.Client{Transport: NewDummyTransport(response)}
	exp := exporter.NewNginxPlusExporter(
		client,
		scraper.NewNginxScraper(),
		scraper.NewNginxPlusScraper(),
		"nginx_test",
		[]string{},
		[]string{"http://localhost:9000"},
	)

	metrics := make(chan prometheus.Metric)

	go func() {
		exp.Collect(metrics)
		close(metrics)
	}()

	checks := map[string]bool{
		"nginx_test_exporter_scrapes_total":                        false,
		"nginx_test_upstream_peer_responses_1xx":                   false,
		"nginx_test_cache_stale_bytes":                             false,
		"nginx_test_stream_upstream_peer_healthchecks_downtime":    false,
		"nginx_test_stream_upstream_peer_healthchecks_downstart":   false,
		"nginx_test_stream_upstream_peer_fails":                    false,
		"nginx_test_zone_responses_4xx":                            false,
		"nginx_test_upstream_queue_max_size":                       false,
		"nginx_test_upstream_peer_responses_3xx":                   false,
		"nginx_test_stream_zone_processing":                        false,
		"nginx_test_cache_miss_responses_written":                  false,
		"nginx_test_cache_miss_bytes_written":                      false,
		"nginx_test_ssl_handshakes":                                false,
		"nginx_test_zone_discarded":                                false,
		"nginx_test_upstream_peer_selected":                        false,
		"nginx_test_cache_revalidated_bytes":                       false,
		"nginx_test_cache_expired_bytes_written":                   false,
		"nginx_test_stream_upstream_peer_healthchecks_fails":       false,
		"nginx_test_upstream_peer_requests":                        false,
		"nginx_test_upstream_peer_healthchecks_unhealthy":          false,
		"nginx_test_cache_miss_responses":                          false,
		"nginx_test_cache_size":                                    false,
		"nginx_test_stream_zone_sent":                              false,
		"nginx_test_stream_upstream_peer_received":                 false,
		"nginx_test_connections_idle":                              false,
		"nginx_test_zone_processing":                               false,
		"nginx_test_upstream_peer_responses_total":                 false,
		"nginx_test_upstream_peer_sent":                            false,
		"nginx_test_upstream_peer_header_time":                     false,
		"nginx_test_cache_responses":                               false,
		"nginx_test_stream_zone_connections":                       false,
		"nginx_test_ssl_session_reuses":                            false,
		"nginx_test_zone_sent":                                     false,
		"nginx_test_upstream_peer_active":                          false,
		"nginx_test_upstream_peer_received":                        false,
		"nginx_test_cache_hit_responses":                           false,
		"nginx_test_cache_expired_bytes":                           false,
		"nginx_test_cache_expired_responses_written":               false,
		"nginx_test_stream_upstream_peer_active":                   false,
		"nginx_test_zone_responses_5xx":                            false,
		"nginx_test_zone_responses_total":                          false,
		"nginx_test_upstream_peer_weight":                          false,
		"nginx_test_upstream_peer_healthchecks_fails":              false,
		"nginx_test_stream_upstream_peer_unavail":                  false,
		"nginx_test_stream_upstream_peer_healthchecks_selected":    false,
		"nginx_test_stream_upstream_peer_first_byte_time":          false,
		"nginx_test_stream_upstream_peer_response_time":            false,
		"nginx_test_upstream_peer_downtime":                        false,
		"nginx_test_upstream_peer_downstart":                       false,
		"nginx_test_cache_responses_written":                       false,
		"nginx_test_connections_active":                            false,
		"nginx_test_zone_responses_1xx":                            false,
		"nginx_test_upstream_peer_responses_2xx":                   false,
		"nginx_test_upstream_peer_healthchecks_checks":             false,
		"nginx_test_cache_updating_responses":                      false,
		"nginx_test_processes_respawned":                           false,
		"nginx_test_zone_responses_2xx":                            false,
		"nginx_test_upstream_keepalive":                            false,
		"nginx_test_upstream_queue_size":                           false,
		"nginx_test_upstream_zombies":                              false,
		"nginx_test_upstream_peer_response_time":                   false,
		"nginx_test_stream_zone_received":                          false,
		"nginx_test_stream_upstream_peer_connect_time":             false,
		"nginx_test_upstream_peer_responses_4xx":                   false,
		"nginx_test_cache_revalidated_responses":                   false,
		"nginx_test_cache_bytes":                                   false,
		"nginx_test_cache_bytes_written":                           false,
		"nginx_test_stream_upstream_peer_sent":                     false,
		"nginx_test_stream_upstream_peer_healthchecks_checks":      false,
		"nginx_test_stream_upstream_peer_healthchecks_unhealthy":   false,
		"nginx_test_zone_received":                                 false,
		"nginx_test_upstream_peer_unavail":                         false,
		"nginx_test_cache_max_size":                                false,
		"nginx_test_stream_upstream_peer_connections":              false,
		"nginx_test_cache_miss_bytes":                              false,
		"nginx_test_stream_upstream_zombies":                       false,
		"nginx_test_stream_upstream_peer_weight":                   false,
		"nginx_test_upstream_queue_overflows":                      false,
		"nginx_test_upstream_peer_fails":                           false,
		"nginx_test_cache_hit_bytes":                               false,
		"nginx_test_upstream_peer_responses_5xx":                   false,
		"nginx_test_cache_stale_responses":                         false,
		"nginx_test_connections_accepted":                          false,
		"nginx_test_connections_dropped":                           false,
		"nginx_test_ssl_handshakes_failed":                         false,
		"nginx_test_zone_requests":                                 false,
		"nginx_test_zone_responses_3xx":                            false,
		"nginx_test_upstream_peer_max_conns":                       false,
		"nginx_test_cache_updating_bytes":                          false,
		"nginx_test_cache_expired_responses":                       false,
		"nginx_test_upstream_peer_backup":                          false,
		"nginx_test_upstream_peer_healthchecks_last_passed":        false,
		"nginx_test_cache_cold":                                    false,
		"nginx_test_stream_upstream_peer_backup":                   false,
		"nginx_test_stream_upstream_peer_healthchecks_last_passed": false,
		"nginx_test_upstream_peer_state":                           false,
		"nginx_test_stream_upstream_peer_state":                    false,
	}

	for m := range metrics {
		switch m.(type) {
		case prometheus.Gauge:
			for metricName := range checks {
				if strings.Contains(m.Desc().String(), metricName) {
					checks[metricName] = true
				}
			}
		}
	}

	for metricName, exists := range checks {
		if !exists {
			c.Errorf("didn't find metric '%s'", metricName)
		}
	}
}

func (s NginxExporterSuite) TestInvalidNginxStatsUrl_Fail(c *C) {
	headers := http.Header{}
	headers.Add("Content-Type", "text/plain")
	response := http.Response{
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       NewDummyBody(nginxStats),
	}

	client := &http.Client{Transport: NewDummyTransport(response)}
	exp := exporter.NewNginxPlusExporter(
		client,
		scraper.NewNginxScraper(),
		scraper.NewNginxPlusScraper(),
		"nginx_test",
		[]string{"invalid nginx stats url"},
		[]string{},
	)

	metrics := make(chan prometheus.Metric, 1)

	go func() {
		exp.Collect(metrics)
		close(metrics)
	}()

	c.Assert(len(metrics), Equals, 0)
}

type DummyTransport struct {
	response http.Response
}

func NewDummyTransport(response http.Response) *DummyTransport {
	return &DummyTransport{response: response}
}

func (transport *DummyTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &transport.response, nil
}

type DummyBody struct {
	*strings.Reader
}

func NewDummyBody(body string) DummyBody {
	return DummyBody{Reader: strings.NewReader(body)}
}

func (body DummyBody) Close() error {
	return nil
}
