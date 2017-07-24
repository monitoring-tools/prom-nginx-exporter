package exporter_test

import (
	"net/http"
	"strings"
	"testing"

	"prom-nginx-exporter/exporter"
	"prom-nginx-exporter/scraper"

	"github.com/prometheus/client_golang/prometheus"
	. "gopkg.in/check.v1"
)

func TestNginxExporter(t *testing.T) { TestingT(t) }

type NginxExporterSuite struct{}

var _ = Suite(&NginxExporterSuite{})

var nginxStats = "Active connections: 2\n server accepts handled requests\n8522429 8522429 8641727\nReading: 0 Writing: 1 Waiting: 3"

func (s NginxExporterSuite) TestScrape_Success(c *C) {
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
