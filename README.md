# Nginx Plus Exporter

Nginx Plus stats exporter for Prometheus.

The Nginx Plus Exporter requests the nginx or nginx plus stats from specified endpoints and exposes them for Prometheus consumption.

## Building and running

For linux:

```
$ make build-linux
./linux_amd64/nginx-plus-exporter listen-address="localhost:9005" --metrics-path="/metrics" --namespace="nginxplus" --nginx-stats-urls="localhost:9002/status" --nginx-stats-urls="localhost:9003/status" --nginx-plus-stats-urls="localhost:9004/status"
```

For darwin:

```
$ make build-darwin
./darwin_amd64/nginx-plus-exporter listen-address="localhost:9005" --metrics-path="/metrics" --namespace="nginxplus" --nginx-stats-urls="localhost:9002/status" --nginx-stats-urls="localhost:9003/status" --nginx-plus-stats-urls="localhost:9004/status"
```

Other useful make commands:

The building application for linux with amd65 architecture:
```
$ make build
```

The running unit tests:
```
$ make test
```

The creating docker image:
```
$ make docker
```

The applying go tool to code:
```
$ make fmt
>> formatting source

$ make lint
>> linting source

$ make imports
>> fixing source imports
```

The running all targets:
```
$ make all
```

It will get all necessary dependencies, run unit tests and build application for linux with amd64 architecture.

### Flags

Name                  | Required | Multiple | Default        | Description
--------------------- | -------- | -------- | -------------- | -----------
listen-address        |    no    |    no    | localhost:9001 | Address on which to expose metrics and web interface.
metrics-path          |    no    |    no    | /metrics       | Path under which to expose metrics.
namespace             |    no    |    no    | nginx          | The namespace of metrics.
nginx-stats-urls      |    yes   |    yes   | -              | An array of Nginx URL to gather stats.
nginx-plus-stats-urls |    yes   |    yes   | -              | An array of Nginx Plus URL to gather stats.

## What's exported?
It exports statistics of standart Nginx module (https://nginx.org/en/docs/http/ngx_http_stub_status_module.html) and Nginx Plus module (http://nginx.org/en/docs/http/ngx_http_status_module.html).
