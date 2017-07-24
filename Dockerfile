FROM golang:latest

ADD . $GOPATH/src/prom-nginx-exporter

RUN mkdir -p $GOPATH/src/github.com/Masterminds/glide && \
    cd $GOPATH/src/github.com/Masterminds/glide && \
    git clone https://github.com/Masterminds/glide.git . && \
    make build && \
    mv glide $GOPATH/bin/glide

RUN cd $GOPATH/src/prom-nginx-exporter && \
    make build && \
    mv $GOPATH/src/prom-nginx-exporter/linux_amd64/prom-nginx-exporter $GOPATH/bin/prom-nginx-exporter

CMD ["prom-nginx-exporter", "--nginx-stats-urls=localhost"]
