package scraper

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/monitoring-tools/prom-nginx-exporter/metric"
)

var (
	// errIncorrectNginxStats describes parse error due to invalid content of stats
	errIncorrectNginxStats = errors.New("incorrect nginx stats")
)

// NginxScraper is the main struct of nginx stats scraper
type NginxScraper struct{}

// NewNginxScraper creates new nginx status scraper
func NewNginxScraper() NginxScraper {
	return NginxScraper{}
}

// Scrape scrapes full information from nginx
func (scr *NginxScraper) Scrape(body io.Reader, metrics chan<- metric.Metric, labels map[string]string) error {
	var reader = bufio.NewReader(body)

	err := scr.scrapeActiveConnections(reader, metrics, labels)
	if err != nil {
		if err == io.EOF {
			err = errIncorrectNginxStats
		}
		return err
	}

	err = scr.scrapeAcceptsHandledRequests(reader, metrics, labels)
	if err != nil {
		if err == io.EOF {
			err = errIncorrectNginxStats
		}
		return err
	}

	return scr.scrapeReadingWritingWaiting(reader, metrics, labels)
}

// scrapeActiveConnections scrapes number of active connections
func (scr *NginxScraper) scrapeActiveConnections(reader *bufio.Reader, metrics chan<- metric.Metric, labels map[string]string) error {
	_, err := reader.ReadString(':')
	if err != nil {
		return err
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	active, err := strconv.ParseUint(strings.TrimSpace(line), 10, 64)
	if err != nil {
		return err
	}

	metrics <- metric.NewMetric("active", active, labels)

	return nil
}

// scrapeAcceptsHandledRequests scrapes number of accepts, handled, requests
func (scr *NginxScraper) scrapeAcceptsHandledRequests(reader *bufio.Reader, metrics chan<- metric.Metric, labels map[string]string) error {
	_, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	data := strings.Fields(line)
	if len(data) != 3 {
		return errors.New("unable to parse server accepts, handled, requests stats")
	}

	accepts, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return err
	}
	metrics <- metric.NewMetric("accepts", accepts, labels)

	handled, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return err
	}
	metrics <- metric.NewMetric("handled", handled, labels)

	requests, err := strconv.ParseUint(data[2], 10, 64)
	if err != nil {
		return err
	}
	metrics <- metric.NewMetric("requests", requests, labels)

	return nil
}

// scrapeReadingWritingWaiting scrapes number of reading, writing, waiting requests
func (scr *NginxScraper) scrapeReadingWritingWaiting(reader *bufio.Reader, metrics chan<- metric.Metric, labels map[string]string) error {
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}

	data := strings.Fields(line)
	if len(data) != 6 {
		return errors.New("unable to parse server reading, writing, waiting stats")
	}

	reading, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return err
	}
	metrics <- metric.NewMetric("reading", reading, labels)

	writing, err := strconv.ParseUint(data[3], 10, 64)
	if err != nil {
		return err
	}
	metrics <- metric.NewMetric("writing", writing, labels)

	waiting, err := strconv.ParseUint(data[5], 10, 64)
	if err != nil {
		return err
	}
	metrics <- metric.NewMetric("waiting", waiting, labels)

	return nil
}
