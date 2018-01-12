FROM quay.io/prometheus/busybox:latest

COPY prom-nginx-exporter /bin/prom-nginx-exporter

CMD ["prom-nginx-exporter", "--nginx-stats-urls=localhost"]
