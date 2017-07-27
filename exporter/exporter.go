package exporter

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/monitoring-tools/prom-nginx-exporter/common"
	"github.com/monitoring-tools/prom-nginx-exporter/metric"
	"github.com/monitoring-tools/prom-nginx-exporter/scraper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	// nginxModule is used to define nginx urls with standard module(ngx_http_stub_status_module)
	nginxModule = "nginx"
	// nginxPlusModule is used to define nginx urls with Plus module(ngx_http_status_module)
	nginxPlusModule = "nginxPlus"
)

// nginxPlusExporter is nginx and nginx plus stats exporter
type nginxPlusExporter struct {
	namespace     string
	nginxUrls     []string
	nginxPlusUrls []string

	client           *http.Client
	nginxScraper     scraper.NginxScraper
	nginxPlusScraper scraper.NginxPlusScraper

	duration     prometheus.Summary
	totalScrapes prometheus.Counter
	metrics      map[string]*prometheus.GaugeVec
}

// NewNginxPlusExporter creates nginx and nginx plus stats exporter
func NewNginxPlusExporter(
	client *http.Client,
	nginxScraper scraper.NginxScraper,
	nginxPlusScraper scraper.NginxPlusScraper,
	namespace string,
	nginxUrls []string,
	nginxPlusUrls []string,
) *nginxPlusExporter {

	duration := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: namespace,
		Name:      "last_scrape_duration_seconds",
		Help:      "The last scrape duration.",
	})

	totalScrapes := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "exporter_scrapes_total",
		Help:      "Current total nginx scrapes.",
	})

	return &nginxPlusExporter{
		client:           client,
		namespace:        namespace,
		nginxUrls:        nginxUrls,
		nginxPlusUrls:    nginxPlusUrls,
		nginxScraper:     nginxScraper,
		nginxPlusScraper: nginxPlusScraper,
		duration:         duration,
		totalScrapes:     totalScrapes,
		metrics:          map[string]*prometheus.GaugeVec{},
	}
}

// Describe describes nginx and nginx plus metrics
func (exp *nginxPlusExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, item := range exp.metrics {
		item.Describe(ch)
	}

	ch <- exp.duration.Desc()
	ch <- exp.totalScrapes.Desc()
}

// Collect collects nginx and nginx plus metrics
func (exp *nginxPlusExporter) Collect(ch chan<- prometheus.Metric) {
	exp.save(exp.scrape())
	exp.expose(ch)
}

// scrape scrapes nginx or nginx plus stats for the passed urls
func (exp *nginxPlusExporter) scrape() chan metric.Metric {
	metrics := make(chan metric.Metric)

	go func() {
		now := time.Now().UnixNano()
		exp.totalScrapes.Inc()

		exp.scrapeModule(nginxModule, exp.nginxUrls, metrics)
		exp.scrapeModule(nginxPlusModule, exp.nginxPlusUrls, metrics)

		exp.duration.Observe(float64(time.Now().UnixNano()-now) / 1000000000)

		close(metrics)
	}()

	return metrics
}

// save saves metrics in the internal struct
func (exp *nginxPlusExporter) save(metrics <-chan metric.Metric) {
	for item := range metrics {
		metricKey := exp.namespace + "_" + item.Name

		if _, ok := exp.metrics[metricKey]; !ok {
			gaugeOpt := prometheus.GaugeOpts{
				Namespace: exp.namespace,
				Name:      item.Name,
			}

			labelNames := make([]string, 0, len(item.Labels))
			for labelName := range item.Labels {
				labelNames = append(labelNames, labelName)
			}

			exp.metrics[metricKey] = prometheus.NewGaugeVec(gaugeOpt, labelNames)
		}

		if val, err := common.ConvertValueToFloat64(item.Value); err != nil {
			log.Errorf("Convert error for metric '%s': %s", item.Name, err)
			continue
		} else {
			exp.metrics[metricKey].With(item.Labels).Set(val)
		}
	}
}

// expose returns metrics to base metric channel
func (exp *nginxPlusExporter) expose(ch chan<- prometheus.Metric) {
	ch <- exp.duration
	ch <- exp.totalScrapes

	for _, m := range exp.metrics {
		m.Collect(ch)
	}
}

// scrapeModule scrapes stats for module(nginx or nginx plus)
func (exp *nginxPlusExporter) scrapeModule(module string, urls []string, metrics chan<- metric.Metric) {
	for _, u := range urls {
		addr, err := url.Parse(u)
		if err != nil {
			log.Fatalf("Unable to parse address '%s': %s", u, err)
		}

		labels := map[string]string{
			"port":   addr.Port(),
			"server": addr.Hostname(),
		}

		err = exp.scrapeURL(module, addr, metrics, labels)
		if err != nil {
			log.Error(err)
		}
	}
}

// scrapeURL scrapes stats for passed url
func (exp *nginxPlusExporter) scrapeURL(module string, addr *url.URL, metrics chan<- metric.Metric, labels map[string]string) error {
	resp, err := exp.client.Get(addr.String())
	if err != nil {
		return fmt.Errorf("Error making HTTP request to '%s': %s", addr.String(), err)
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return fmt.Errorf("%s returned HTTP status %d", addr.String(), resp.StatusCode)
	}

	contentType := strings.Split(resp.Header.Get("Content-Type"), ";")[0]

	if module == nginxModule {
		err = exp.nginxScraper.Scrape(resp.Body, metrics, labels)
		if err != nil {
			return fmt.Errorf("Error scraping nginx stats using address '%s': %s", addr.String(), err)
		}

		return nil
	} else if module == nginxPlusModule && contentType == "application/json" {
		err = exp.nginxPlusScraper.Scrape(resp.Body, metrics, labels)
		if err != nil {
			return fmt.Errorf("Error scraping nginx plus stats using address '%s': %s", addr.String(), err)
		}

		return nil
	}

	return fmt.Errorf("%s returned unsupported content type '%s'", addr.String(), contentType)
}
