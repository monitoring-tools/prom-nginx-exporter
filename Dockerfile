FROM golang:latest

ADD . $GOPATH/src/prom-nginx-exporter

RUN go get -u github.com/golang/dep/cmd/dep

RUN cd $GOPATH/src/prom-nginx-exporter && \
    make build && \
    mv $GOPATH/src/prom-nginx-exporter/linux_amd64/prom-nginx-exporter $GOPATH/bin/prom-nginx-exporter

CMD ["prom-nginx-exporter", "--nginx-stats-urls=localhost"]
