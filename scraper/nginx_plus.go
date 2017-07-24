package scraper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"prom-nginx-exporter/metric"
)

// NginxPlusScraper is scraper for getting nginx plus metrics
type NginxPlusScraper struct{}

// NewNginxPlusScraper crates new nginx plus stats scraper
func NewNginxPlusScraper() NginxPlusScraper {
	return NginxPlusScraper{}
}

// Scrape scrapes stats from nginx plus module
func (scr *NginxPlusScraper) Scrape(body io.Reader, metrics chan<- metric.Metric, labels map[string]string) error {
	dec := json.NewDecoder(bufio.NewReader(body))

	status := &Status{}
	if err := dec.Decode(status); err != nil {
		return fmt.Errorf("Error while decoding JSON response")
	}

	scr.scrapeProcesses(status, metrics, labels)
	scr.scrapeConnections(status, metrics, labels)
	scr.scrapeSsl(status, metrics, labels)
	scr.scrapeRequest(status, metrics, labels)
	scr.scrapeUpstream(status, metrics, labels)
	scr.scrapeCache(status, metrics, labels)
	scr.scrapeStream(status, metrics, labels)

	return nil
}

// scrapeProcesses scrapes processes metrics
func (scr *NginxPlusScraper) scrapeProcesses(status *Status, metrics chan<- metric.Metric, labels map[string]string) {
	metrics <- metric.NewMetric("processes_respawned", *status.Processes.Respawned, labels)
}

// scrapeConnections scrapes connections metrics
func (scr *NginxPlusScraper) scrapeConnections(status *Status, metrics chan<- metric.Metric, labels map[string]string) {
	metrics <- metric.NewMetric("connections_accepted", status.Connections.Accepted, labels)
	metrics <- metric.NewMetric("connections_dropped", status.Connections.Dropped, labels)
	metrics <- metric.NewMetric("connections_active", status.Connections.Active, labels)
	metrics <- metric.NewMetric("connections_idle", status.Connections.Idle, labels)
}

// scrapeSsl scrapes SSL metrics
func (scr *NginxPlusScraper) scrapeSsl(status *Status, metrics chan<- metric.Metric, labels map[string]string) {
	metrics <- metric.NewMetric("ssl_handshakes", status.Ssl.Handshakes, labels)
	metrics <- metric.NewMetric("ssl_handshakes_failed", status.Ssl.HandshakesFailed, labels)
	metrics <- metric.NewMetric("ssl_session_reuses", status.Ssl.SessionReuses, labels)
}

// scrapeRequest scrapes request metrics
func (scr *NginxPlusScraper) scrapeRequest(status *Status, metrics chan<- metric.Metric, labels map[string]string) {
	for zoneName, zone := range status.ServerZones {
		zoneLabels := make(map[string]string)
		for k, v := range labels {
			zoneLabels[k] = v
		}

		zoneLabels["zone"] = zoneName
		metrics <- metric.NewMetric("zone_processing", zone.Processing, zoneLabels)
		metrics <- metric.NewMetric("zone_requests", zone.Requests, zoneLabels)
		metrics <- metric.NewMetric("zone_responses_1xx", zone.Responses.Responses1xx, zoneLabels)
		metrics <- metric.NewMetric("zone_responses_2xx", zone.Responses.Responses2xx, zoneLabels)
		metrics <- metric.NewMetric("zone_responses_3xx", zone.Responses.Responses3xx, zoneLabels)
		metrics <- metric.NewMetric("zone_responses_4xx", zone.Responses.Responses4xx, zoneLabels)
		metrics <- metric.NewMetric("zone_responses_5xx", zone.Responses.Responses5xx, zoneLabels)
		metrics <- metric.NewMetric("zone_responses_total", zone.Responses.Total, zoneLabels)
		metrics <- metric.NewMetric("zone_received", zone.Received, zoneLabels)
		metrics <- metric.NewMetric("zone_sent", zone.Sent, zoneLabels)

		if zone.Discarded != nil {
			metrics <- metric.NewMetric("zone_discarded", *zone.Discarded, zoneLabels)
		}
	}
}

// scrapeUpstream scrapes upstream metrics
func (scr *NginxPlusScraper) scrapeUpstream(status *Status, metrics chan<- metric.Metric, labels map[string]string) {
	for upstreamName, upstream := range status.Upstreams {
		upstreamLabels := make(map[string]string)
		for k, v := range labels {
			upstreamLabels[k] = v
		}

		upstreamLabels["upstream"] = upstreamName

		metrics <- metric.NewMetric("upstream_keepalive", upstream.Keepalive, upstreamLabels)
		metrics <- metric.NewMetric("upstream_zombies", upstream.Zombies, upstreamLabels)

		if upstream.Queue != nil {
			metrics <- metric.NewMetric("upstream_queue_size", upstream.Queue.Size, upstreamLabels)
			metrics <- metric.NewMetric("upstream_queue_max_size", upstream.Queue.MaxSize, upstreamLabels)
			metrics <- metric.NewMetric("upstream_queue_overflows", upstream.Queue.Overflows, upstreamLabels)
		}

		for _, peer := range upstream.Peers {
			peerLabels := make(map[string]string)
			for k, v := range upstreamLabels {
				peerLabels[k] = v
			}
			peerLabels["serverAddress"] = peer.Server
			if peer.ID != nil {
				peerLabels["id"] = strconv.Itoa(*peer.ID)
			}

			metrics <- metric.NewMetric("upstream_peer_backup", peer.Backup, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_weight", peer.Weight, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_state", peer.State, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_active", peer.Active, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_requests", peer.Requests, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_responses_1xx", peer.Responses.Responses1xx, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_responses_2xx", peer.Responses.Responses2xx, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_responses_3xx", peer.Responses.Responses3xx, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_responses_4xx", peer.Responses.Responses4xx, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_responses_5xx", peer.Responses.Responses5xx, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_responses_total", peer.Responses.Total, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_sent", peer.Sent, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_received", peer.Received, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_fails", peer.Fails, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_unavail", peer.Unavail, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_healthchecks_checks", peer.HealthChecks.Checks, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_healthchecks_fails", peer.HealthChecks.Fails, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_healthchecks_unhealthy", peer.HealthChecks.Unhealthy, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_downtime", peer.Downtime, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_downstart", peer.Downstart, peerLabels)
			metrics <- metric.NewMetric("upstream_peer_selected", *peer.Selected, peerLabels)

			if peer.HealthChecks.LastPassed != nil {
				metrics <- metric.NewMetric("upstream_peer_healthchecks_last_passed", *peer.HealthChecks.LastPassed, peerLabels)
			}

			if peer.HeaderTime != nil {
				metrics <- metric.NewMetric("upstream_peer_header_time", *peer.HeaderTime, peerLabels)
			}

			if peer.ResponseTime != nil {
				metrics <- metric.NewMetric("upstream_peer_response_time", *peer.ResponseTime, peerLabels)
			}

			if peer.MaxConns != nil {
				metrics <- metric.NewMetric("upstream_peer_max_conns", *peer.MaxConns, peerLabels)
			}
		}
	}
}

// scrapeCache scrapes cache metrics
func (scr *NginxPlusScraper) scrapeCache(status *Status, metrics chan<- metric.Metric, labels map[string]string) {
	for cacheName, cache := range status.Caches {
		cacheLabels := make(map[string]string)
		for k, v := range labels {
			cacheLabels[k] = v
		}
		cacheLabels["cache"] = cacheName

		metrics <- metric.NewMetric("cache_size", cache.Size, cacheLabels)
		metrics <- metric.NewMetric("cache_max_size", cache.Size, cacheLabels)
		metrics <- metric.NewMetric("cache_cold", cache.Cold, cacheLabels)
		metrics <- metric.NewMetric("cache_hit_responses", cache.Hit.Responses, cacheLabels)
		metrics <- metric.NewMetric("cache_hit_bytes", cache.Hit.Bytes, cacheLabels)
		metrics <- metric.NewMetric("cache_stale_responses", cache.Stale.Responses, cacheLabels)
		metrics <- metric.NewMetric("cache_stale_bytes", cache.Stale.Bytes, cacheLabels)
		metrics <- metric.NewMetric("cache_updating_responses", cache.Updating.Responses, cacheLabels)
		metrics <- metric.NewMetric("cache_updating_bytes", cache.Updating.Bytes, cacheLabels)
		metrics <- metric.NewMetric("cache_revalidated_responses", cache.Revalidated.Responses, cacheLabels)
		metrics <- metric.NewMetric("cache_revalidated_bytes", cache.Revalidated.Bytes, cacheLabels)
		metrics <- metric.NewMetric("cache_miss_responses", cache.Miss.Responses, cacheLabels)
		metrics <- metric.NewMetric("cache_miss_bytes", cache.Miss.Bytes, cacheLabels)
		metrics <- metric.NewMetric("cache_miss_responses_written", cache.Miss.ResponsesWritten, cacheLabels)
		metrics <- metric.NewMetric("cache_miss_bytes_written", cache.Miss.Bytes, cacheLabels)
		metrics <- metric.NewMetric("cache_expired_responses", cache.Expired.Responses, cacheLabels)
		metrics <- metric.NewMetric("cache_expired_bytes", cache.Expired.Bytes, cacheLabels)
		metrics <- metric.NewMetric("cache_expired_responses_written", cache.Expired.ResponsesWritten, cacheLabels)
		metrics <- metric.NewMetric("cache_expired_bytes_written", cache.Expired.BytesWritten, cacheLabels)
		metrics <- metric.NewMetric("cache_responses", cache.Bypass.Responses, cacheLabels)
		metrics <- metric.NewMetric("cache_bytes", cache.Bypass.Bytes, cacheLabels)
		metrics <- metric.NewMetric("cache_responses_written", cache.Bypass.ResponsesWritten, cacheLabels)
		metrics <- metric.NewMetric("cache_bytes_written", cache.Bypass.BytesWritten, cacheLabels)
	}
}

// scrapeStream scrapes stream metrics
func (scr *NginxPlusScraper) scrapeStream(status *Status, metrics chan<- metric.Metric, labels map[string]string) {
	for zoneName, zone := range status.Stream.ServerZones {
		zoneLabels := map[string]string{}
		for k, v := range labels {
			zoneLabels[k] = v
		}
		zoneLabels["zone"] = zoneName

		metrics <- metric.NewMetric("stream_zone_processing", zone.Processing, zoneLabels)
		metrics <- metric.NewMetric("stream_zone_connections", zone.Connections, zoneLabels)
		metrics <- metric.NewMetric("stream_zone_received", zone.Received, zoneLabels)
		metrics <- metric.NewMetric("stream_zone_sent", zone.Sent, zoneLabels)
	}

	for upstreamName, upstream := range status.Stream.Upstreams {
		upstreamLabels := map[string]string{}
		for k, v := range labels {
			upstreamLabels[k] = v
		}
		upstreamLabels["upstream"] = upstreamName

		metrics <- metric.NewMetric("stream_upstream_zombies", upstream.Zombies, upstreamLabels)

		for _, peer := range upstream.Peers {
			peerLables := map[string]string{}
			for k, v := range upstreamLabels {
				peerLables[k] = v
			}
			peerLables["serverAddress"] = peer.Server
			peerLables["id"] = strconv.Itoa(peer.ID)

			metrics <- metric.NewMetric("stream_upstream_peer_backup", peer.Backup, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_weight", peer.Weight, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_state", peer.State, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_active", peer.Active, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_connections", peer.Connections, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_sent", peer.Sent, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_received", peer.Received, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_fails", peer.Fails, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_unavail", peer.Unavail, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_healthchecks_checks", peer.HealthChecks.Checks, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_healthchecks_fails", peer.HealthChecks.Fails, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_healthchecks_unhealthy", peer.HealthChecks.Unhealthy, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_healthchecks_downtime", peer.Downtime, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_healthchecks_downstart", peer.Downstart, peerLables)
			metrics <- metric.NewMetric("stream_upstream_peer_healthchecks_selected", peer.Selected, peerLables)

			if peer.HealthChecks.LastPassed != nil {
				metrics <- metric.NewMetric("stream_upstream_peer_healthchecks_last_passed", *peer.HealthChecks.LastPassed, peerLables)
			}
			if peer.ConnectTime != nil {
				metrics <- metric.NewMetric("stream_upstream_peer_connect_time", *peer.ConnectTime, peerLables)
			}
			if peer.FirstByteTime != nil {
				metrics <- metric.NewMetric("stream_upstream_peer_first_byte_time", *peer.FirstByteTime, peerLables)
			}
			if peer.ResponseTime != nil {
				metrics <- metric.NewMetric("stream_upstream_peer_response_time", *peer.ResponseTime, peerLables)
			}
		}
	}
}

/*
Structures built based on history of status module documentation
http://nginx.org/en/docs/http/ngx_http_status_module.html
Subsequent versions of status response structure available here:
1. http://web.archive.org/web/20130805111222/http://nginx.org/en/docs/http/ngx_http_status_module.html
2. http://web.archive.org/web/20131218101504/http://nginx.org/en/docs/http/ngx_http_status_module.html
3. not available
4. http://web.archive.org/web/20141218170938/http://nginx.org/en/docs/http/ngx_http_status_module.html
5. http://web.archive.org/web/20150414043916/http://nginx.org/en/docs/http/ngx_http_status_module.html
6. http://web.archive.org/web/20150918163811/http://nginx.org/en/docs/http/ngx_http_status_module.html
7. http://web.archive.org/web/20161107221028/http://nginx.org/en/docs/http/ngx_http_status_module.html
*/

// Status is the main struct of nginx plus statistics.
type Status struct {
	Version       int    `json:"version"`
	NginxVersion  string `json:"nginx_version"`
	Address       string `json:"address"`
	Generation    *int   `json:"generation"`     // added in version 5
	LoadTimestamp *int64 `json:"load_timestamp"` // added in version 2
	Timestamp     int64  `json:"timestamp"`
	Pid           *int   `json:"pid"` // added in version 6

	Processes   *Processes  `json:"processes"`
	Connections Connections `json:"connections"`
	Ssl         *Ssl        `json:"ssl"`
	Requests    Requests    `json:"requests"`
	ServerZones ServerZones `json:"server_zones"`
	Upstreams   Upstreams   `json:"upstreams"`
	Caches      Caches      `json:"caches"`
	Stream      Stream      `json:"stream"`
}

// Processes contains the total number of respawned child processes.
type Processes struct {
	// added in version 5
	Respawned *int `json:"respawned"`
}

// Connections contains the total number of accepted, dropped, active and idle client connections.
type Connections struct {
	Accepted int `json:"accepted"`
	Dropped  int `json:"dropped"`
	Active   int `json:"active"`
	Idle     int `json:"idle"`
}

// Ssl contains the total number of successful, failed SSL handshakes and number of sessions reuses during SSL handshake.
type Ssl struct {
	// added in version 6
	Handshakes       int64 `json:"handshakes"`
	HandshakesFailed int64 `json:"handshakes_failed"`
	SessionReuses    int64 `json:"session_reuses"`
}

// Requests contains total and current number of client requests.
type Requests struct {
	Total   int64 `json:"total"`
	Current int   `json:"current"`
}

// ServerZones contains info about processed requests, requests received from clients, number of responses from clients
// with http statuses, total number of requests completed without sending a response, number of bytes received and sent.
type ServerZones map[string]struct {
	// added in version 2
	Processing int   `json:"processing"`
	Requests   int64 `json:"requests"`
	Responses  struct {
		Responses1xx int64 `json:"1xx"`
		Responses2xx int64 `json:"2xx"`
		Responses3xx int64 `json:"3xx"`
		Responses4xx int64 `json:"4xx"`
		Responses5xx int64 `json:"5xx"`
		Total        int64 `json:"total"`
	} `json:"responses"`
	Discarded *int64 `json:"discarded"` // added in version 6
	Received  int64  `json:"received"`
	Sent      int64  `json:"sent"`
}

// Upstreams contains a lot of information about upstreams, like: peers info, current number of idle keepalive
// connections, total number of zombies, the size of requests queue.
type Upstreams map[string]struct {
	Peers []struct {
		ID        *int   `json:"id"` // added in version 3
		Server    string `json:"server"`
		Backup    bool   `json:"backup"`
		Weight    int    `json:"weight"`
		State     string `json:"state"`
		Active    int    `json:"active"`
		Keepalive *int   `json:"keepalive"` // removed in version 5
		MaxConns  *int   `json:"max_conns"` // added in version 3
		Requests  int64  `json:"requests"`
		Responses struct {
			Responses1xx int64 `json:"1xx"`
			Responses2xx int64 `json:"2xx"`
			Responses3xx int64 `json:"3xx"`
			Responses4xx int64 `json:"4xx"`
			Responses5xx int64 `json:"5xx"`
			Total        int64 `json:"total"`
		} `json:"responses"`
		Sent         int64 `json:"sent"`
		Received     int64 `json:"received"`
		Fails        int64 `json:"fails"`
		Unavail      int64 `json:"unavail"`
		HealthChecks struct {
			Checks     int64 `json:"checks"`
			Fails      int64 `json:"fails"`
			Unhealthy  int64 `json:"unhealthy"`
			LastPassed *bool `json:"last_passed"`
		} `json:"health_checks"`
		Downtime     int64  `json:"downtime"`
		Downstart    int64  `json:"downstart"`
		Selected     *int64 `json:"selected"`      // added in version 4
		HeaderTime   *int64 `json:"header_time"`   // added in version 5
		ResponseTime *int64 `json:"response_time"` // added in version 5
	} `json:"peers"`
	Keepalive int `json:"keepalive"`
	Zombies   int `json:"zombies"` // added in version 6
	Queue     *struct {
		// added in version 6
		Size      int   `json:"size"`
		MaxSize   int   `json:"max_size"`
		Overflows int64 `json:"overflows"`
	} `json:"queue"`
}

// Caches contains a lot of information of cache, like: current size of cache, the limit on the maximum size of the
// cache, total number of responses and total number of bytes read from the cache, total number of requests not read
// from the cache and number of bytes read from proxied server, number of responses and bytes written to the cache.
type Caches map[string]struct {
	// added in version 2
	Size    int64 `json:"size"`
	MaxSize int64 `json:"max_size"`
	Cold    bool  `json:"cold"`
	Hit     struct {
		Responses int64 `json:"responses"`
		Bytes     int64 `json:"bytes"`
	} `json:"hit"`
	Stale struct {
		Responses int64 `json:"responses"`
		Bytes     int64 `json:"bytes"`
	} `json:"stale"`
	Updating struct {
		Responses int64 `json:"responses"`
		Bytes     int64 `json:"bytes"`
	} `json:"updating"`
	Revalidated *struct {
		// added in version 3
		Responses int64 `json:"responses"`
		Bytes     int64 `json:"bytes"`
	} `json:"revalidated"`
	Miss struct {
		Responses        int64 `json:"responses"`
		Bytes            int64 `json:"bytes"`
		ResponsesWritten int64 `json:"responses_written"`
		BytesWritten     int64 `json:"bytes_written"`
	} `json:"miss"`
	Expired struct {
		Responses        int64 `json:"responses"`
		Bytes            int64 `json:"bytes"`
		ResponsesWritten int64 `json:"responses_written"`
		BytesWritten     int64 `json:"bytes_written"`
	} `json:"expired"`
	Bypass struct {
		Responses        int64 `json:"responses"`
		Bytes            int64 `json:"bytes"`
		ResponsesWritten int64 `json:"responses_written"`
		BytesWritten     int64 `json:"bytes_written"`
	} `json:"bypass"`
}

// Stream contains a lot of information about streams, like: the number of processed client connections, number of
// accepted connections, number of completed client sessions with http statuses, number of connections without creating
// a session, number of bytes received from clients, total number of bytes sent to clients and more information about
// upstreams.
type Stream struct {
	ServerZones map[string]struct {
		Processing  int `json:"processing"`
		Connections int `json:"connections"`
		Sessions    *struct {
			Total       int64 `json:"total"`
			Sessions1xx int64 `json:"1xx"`
			Sessions2xx int64 `json:"2xx"`
			Sessions3xx int64 `json:"3xx"`
			Sessions4xx int64 `json:"4xx"`
			Sessions5xx int64 `json:"5xx"`
		} `json:"sessions"`
		Discarded *int64 `json:"discarded"` // added in version 7
		Received  int64  `json:"received"`
		Sent      int64  `json:"sent"`
	} `json:"server_zones"`
	Upstreams map[string]struct {
		Peers []struct {
			ID            int    `json:"id"`
			Server        string `json:"server"`
			Backup        bool   `json:"backup"`
			Weight        int    `json:"weight"`
			State         string `json:"state"`
			Active        int    `json:"active"`
			Connections   int64  `json:"connections"`
			ConnectTime   *int   `json:"connect_time"`
			FirstByteTime *int   `json:"first_byte_time"`
			ResponseTime  *int   `json:"response_time"`
			Sent          int64  `json:"sent"`
			Received      int64  `json:"received"`
			Fails         int64  `json:"fails"`
			Unavail       int64  `json:"unavail"`
			HealthChecks  struct {
				Checks     int64 `json:"checks"`
				Fails      int64 `json:"fails"`
				Unhealthy  int64 `json:"unhealthy"`
				LastPassed *bool `json:"last_passed"`
			} `json:"health_checks"`
			Downtime  int64 `json:"downtime"`
			Downstart int64 `json:"downstart"`
			Selected  int64 `json:"selected"`
		} `json:"peers"`
		Zombies int `json:"zombies"`
	} `json:"upstreams"`
}
