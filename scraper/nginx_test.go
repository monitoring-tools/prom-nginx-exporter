package scraper_test

import (
	"strings"
	"testing"

	"github.com/monitoring-tools/prom-nginx-exporter/common"
	"github.com/monitoring-tools/prom-nginx-exporter/metric"
	"github.com/monitoring-tools/prom-nginx-exporter/scraper"
	. "gopkg.in/check.v1"
)

func TestNginScraper(t *testing.T) { TestingT(t) }

type NginxScraperSuite struct{}

var _ = Suite(&NginxScraperSuite{})

func (s NginxScraperSuite) TestScrape_Success(c *C) {
	nginxScrape := scraper.NewNginxScraper()
	reader := strings.NewReader("Active connections: 2\n" +
		"server accepts handled requests\n" +
		"8522429 8522429 8641727\n" +
		"Reading: 0 Writing: 1 Waiting: 3")

	metrics := make(chan metric.Metric, 7)
	labels := map[string]string{
		"host": "localhost",
		"port": "8080",
	}

	err := nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, IsNil, Commentf("error occurred during scrape nginx stats"))

	// assert number of active connections
	m := <-metrics
	c.Assert(m.Name, Equals, "active", Commentf("incorrect name of active connections"))
	active, err := common.ConvertValueToFloat64(m.Value)
	c.Assert(err, IsNil, Commentf("error occurred during convert 'active' metric value to float64"))
	c.Assert(active, Equals, float64(2), Commentf("incorrect number of active connections"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	// assert number of accept connections
	m = <-metrics
	c.Assert(m.Name, Equals, "accepts", Commentf("incorrect name of accepts connections"))
	accepts, err := common.ConvertValueToFloat64(m.Value)
	c.Assert(err, IsNil, Commentf("error occurred during convert 'accepts' metric value to float64"))
	c.Assert(accepts, Equals, float64(8522429), Commentf("incorrect number of accepts connections"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	// assert number of handled connections
	m = <-metrics
	c.Assert(m.Name, Equals, "handled", Commentf("incorrect name of handled connections"))
	handled, err := common.ConvertValueToFloat64(m.Value)
	c.Assert(err, IsNil, Commentf("error occurred during convert 'handled' metric value to float64"))
	c.Assert(handled, Equals, float64(8522429), Commentf("incorrect number of handled connections"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	// assert number of requests connections
	m = <-metrics
	c.Assert(m.Name, Equals, "requests", Commentf("incorrect name of requests connections"))
	requests, err := common.ConvertValueToFloat64(m.Value)
	c.Assert(err, IsNil, Commentf("error occurred during convert 'requests' metric value to float64"))
	c.Assert(requests, Equals, float64(8641727), Commentf("incorrect number of requests connections"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	// assert number of reading connections
	m = <-metrics
	c.Assert(m.Name, Equals, "reading", Commentf("incorrect name of reading connections"))
	reading, err := common.ConvertValueToFloat64(m.Value)
	c.Assert(err, IsNil, Commentf("error occurred during convert 'reading' metric value to float64"))
	c.Assert(reading, Equals, float64(0), Commentf("incorrect number of reading connections"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	// assert number of writing connections
	m = <-metrics
	c.Assert(m.Name, Equals, "writing", Commentf("incorrect name of writing connections"))
	writing, err := common.ConvertValueToFloat64(m.Value)
	c.Assert(err, IsNil, Commentf("error occurred during convert 'writing' metric value to float64"))
	c.Assert(writing, Equals, float64(1), Commentf("incorrect number of writing connections"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))

	// assert number of waiting connections
	m = <-metrics
	c.Assert(m.Name, Equals, "waiting", Commentf("incorrect name of waiting connections"))
	waiting, err := common.ConvertValueToFloat64(m.Value)
	c.Assert(err, IsNil, Commentf("error occurred during convert 'waiting' metric value to float64"))
	c.Assert(waiting, Equals, float64(3), Commentf("incorrect number of waiting connections"))
	c.Assert(m.Labels, DeepEquals, labels, Commentf("incorrect set of labels"))
}

func (s NginxScraperSuite) TestScrape_ActiveConnections_Fail(c *C) {
	nginxScrape := scraper.NewNginxScraper()
	metrics := make(chan metric.Metric, 0)
	labels := make(map[string]string)

	reader := strings.NewReader("Active connections 2\n" +
		"server accepts handled requests\n" +
		"8522429 8522429 8641727\n" +
		"Reading 0 Writing 1 Waiting 3")

	err := nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("error occurred during parse active connections stats"))
	c.Assert(err.Error(), Equals, "incorrect nginx stats", Commentf("common is skipped in active connection line"))

	metrics = make(chan metric.Metric, 1)
	reader = strings.NewReader("Active connections: 2 server accepts handled requests 8522429 8522429 8641727 Reading: 0 Writing: 1 Waiting: 3")
	err = nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("error occurred during parse active connections stat"))
	c.Assert(err.Error(), Equals, "incorrect nginx stats", Commentf("there is no new line after active connection info"))

	metrics = make(chan metric.Metric, 1)
	reader = strings.NewReader("Active connections: str\n" +
		"server accepts handled requests\n" +
		"8522429 8522429 8641727\n" +
		"Reading: 0 Writing: 1 Waiting: 3")
	err = nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("should be error of parsing active connections"))
	c.Assert(err.Error(), Equals, "strconv.ParseUint: parsing \"str\": invalid syntax", Commentf("error occurred during parse active connections"))
}

func (s NginxScraperSuite) TestScrapeAcceptsHandledRequests_Fail(c *C) {
	nginxScrape := scraper.NewNginxScraper()
	metrics := make(chan metric.Metric, 1)
	labels := make(map[string]string)
	reader := strings.NewReader("Active connections: 2\n" +
		"server accepts handled requests 8522429 8522429 8641727 Reading: 0 Writing: 1 Waiting: 3")
	err := nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("should be error of parsing header of accepts handled requests"))
	c.Assert(err.Error(), Equals, "incorrect nginx stats", Commentf("error occurred during parse header of accepts handled requests"))

	metrics = make(chan metric.Metric, 1)
	reader = strings.NewReader("Active connections: 2\n" +
		"server accepts handled requests\n" +
		"accepts_str 8522429 8641727\n" +
		"Reading: 0 Writing: 1 Waiting: 3")
	err = nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("should be error of parsing acceps"))
	c.Assert(err.Error(), Equals, "strconv.ParseUint: parsing \"accepts_str\": invalid syntax", Commentf("error occurred during parse accepts"))

	metrics = make(chan metric.Metric, 2)
	reader = strings.NewReader("Active connections: 2\n" +
		"server accepts handled requests\n" +
		"8522429 handled_str 8641727\n" +
		"Reading: 0 Writing: 1 Waiting: 3")
	err = nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("should be error of parsing handled"))
	c.Assert(err.Error(), Equals, "strconv.ParseUint: parsing \"handled_str\": invalid syntax", Commentf("error occurred during parse handled"))

	metrics = make(chan metric.Metric, 3)
	reader = strings.NewReader("Active connections: 2\n" +
		"server accepts handled requests\n" +
		"8522429 8522429 requests_str\n" +
		"Reading: 0 Writing: 1 Waiting: 3")
	err = nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("should be error of parsing requests"))
	c.Assert(err.Error(), Equals, "strconv.ParseUint: parsing \"requests_str\": invalid syntax", Commentf("error occurred during parse requests"))
}

func (s NginxScraperSuite) TestScrapeReadingWritingWaiting_Fail(c *C) {
	nginxScrape := scraper.NewNginxScraper()
	metrics := make(chan metric.Metric, 4)
	labels := make(map[string]string)
	reader := strings.NewReader("Active connections: 2\n" +
		"server accepts handled requests\n" +
		"8522429 8522429 8641727\n" +
		"Reading: reading_str Writing: 1 Waiting: 3")
	err := nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("should be error of parsing reading"))
	c.Assert(err.Error(), Equals, "strconv.ParseUint: parsing \"reading_str\": invalid syntax", Commentf("error occurred during parse reading"))

	metrics = make(chan metric.Metric, 5)
	reader = strings.NewReader("Active connections: 2\n" +
		"server accepts handled requests\n" +
		"8522429 8522429 8641727\n" +
		"Reading: 0 Writing: writing_str Waiting: 3")
	err = nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("should be error of parsing writing"))
	c.Assert(err.Error(), Equals, "strconv.ParseUint: parsing \"writing_str\": invalid syntax", Commentf("error occurred during parse writing"))

	metrics = make(chan metric.Metric, 6)
	reader = strings.NewReader("Active connections: 2\n" +
		"server accepts handled requests\n" +
		"8522429 8522429 8641727\n" +
		"Reading: 0 Writing: 1 Waiting: waiting_str")
	err = nginxScrape.Scrape(reader, metrics, labels)
	c.Assert(err, NotNil, Commentf("should be error of parsing waiting"))
	c.Assert(err.Error(), Equals, "strconv.ParseUint: parsing \"waiting_str\": invalid syntax", Commentf("error occurred during parse waiting"))
}
