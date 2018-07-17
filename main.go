package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/monitoring-tools/prom-nginx-exporter/common"
	"github.com/monitoring-tools/prom-nginx-exporter/exporter"
	"github.com/monitoring-tools/prom-nginx-exporter/scraper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (
	gitSummary string
)

var (
	landingPage = `<html>
<head>
<title>Prom Nginx exporter</title>
</head>
<body>
<h1>Prom Nginx exporter</h1>
<p><a href="/metrics">Metrics</a></p>
</body>
</html>`
)

func main() {
	var config, err = parseFlag()
	if err != nil {
		log.Fatalln(err)
	}

	registerExporter(config.Namespace, config.NginxUrls, config.NginxPlusUrls, config.ExcludeUpstreamPeers)
	run(config.ListenAddress, config.MetricsPath)
}

// parseFlag parses config parameters
func parseFlag() (*common.Config, error) {
	var (
		listenAddress        *string
		metricsPath          *string
		namespace            *string
		version              *bool
		nginxUrls            common.ArrFlags
		nginxPlusUrls        common.ArrFlags
		excludeUpstreamPeers common.ArrFlags
	)

	listenAddress = flag.String("listen-address", ":9001", "Address on which to expose metrics and web interface.")
	metricsPath = flag.String("metrics-path", "/metrics", "Path under which to expose metrics.")
	namespace = flag.String("namespace", "nginx", "The namespace of metrics.")
	version = flag.Bool("version", false, "The version of the exporter.")
	flag.Var(&nginxUrls, "nginx-stats-urls", "An array of Nginx status URLs to gather stats.")
	flag.Var(&nginxPlusUrls, "nginx-plus-stats-urls", "An array of Nginx Plus status URLs to gather stats.")
	flag.Var(&excludeUpstreamPeers, "exclude-upstream-peers", "An array of upstream addresses that need to be excluded in gathering stats.")

	flag.Parse()

	if *version {
		fmt.Println(gitSummary)
		os.Exit(0)
	}

	if len(nginxUrls) == 0 && len(nginxPlusUrls) == 0 {
		return nil, errors.New("no nginx or nginx plus stats url specified")
	}

	return common.NewConfig(*listenAddress, *metricsPath, *namespace, nginxUrls, nginxPlusUrls, excludeUpstreamPeers), nil
}

// registerExporter registers custom nginx metrics exporter
func registerExporter(namespace string, nginxUrls []string, nginxPlusUrls []string, excludeUpstreamPeers []string) {
	var (
		transport = &http.Transport{ResponseHeaderTimeout: time.Duration(3 * time.Second)}
		client    = &http.Client{Transport: transport, Timeout: time.Duration(4 * time.Second)}
	)

	prometheus.MustRegister(exporter.NewNginxPlusExporter(
		client,
		scraper.NewNginxScraper(),
		scraper.NewNginxPlusScraper(excludeUpstreamPeers),
		namespace,
		nginxUrls,
		nginxPlusUrls,
	))
}

// run runs exporter
func run(listenAddress, metricsPath string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(landingPage))
	})
	http.Handle(metricsPath, promhttp.Handler())

	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
