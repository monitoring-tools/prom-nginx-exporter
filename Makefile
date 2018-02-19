DESTDIR ?=
PREFIX ?= /usr/local

GO              = go
GIT_SUMMARY		:= $(shell git describe --tags --always 2>/dev/null)
CURRENT_GOPATH	:=$(lastword $(subst :, ,$(GOPATH)))
DEP_BIN			:= $(GOPATH)/bin/dep

all: deps test build

deps: $(DEP_BIN) Gopkg.toml
	@echo ">> installing application dependencies"
	$(DEP_BIN) ensure

build: linux_amd64/prom-nginx-exporter

build-linux: linux_amd64/prom-nginx-exporter

build-darwin: darwin_amd64/prom-nginx-exporter

linux_amd64/prom-nginx-exporter: deps
	@echo ">> building $@ with version $(GIT_SUMMARY)"
	@GOOS="linux" GOARCH="amd64" $(GO) build -ldflags="-X 'main.gitSummary=$(GIT_SUMMARY)'" -o $@ main.go

darwin_amd64/prom-nginx-exporter: deps
	@echo ">> building $@ with version $(GIT_SUMMARY)"
	@GOOS="darwin" GOARCH="amd64" $(GO) build -ldflags="-X 'main.gitSummary=$(GIT_SUMMARY)'" -o $@ main.go

test:
	@echo ">> making tests"
	@$(GO) test $$($(GO) list ./... | grep -v /vendor/)

docker: Dockerfile
	@echo ">> building docker image prom-nginx-exporter:$(GIT_SUMMARY)"
	@docker build -t prom-nginx-exporter:$(GIT_SUMMARY) .

install: build
	install -m 0755 -d $(DESTDIR)$(PREFIX)/bin
	install -m 0755 linux_amd64/prom-nginx-exporter $(DESTDIR)$(PREFIX)/bin

clean:
	@echo ">> removing build directory"
	@rm -rf vendor/ linux_amd64/ darwin_amd64/

fmt:
	@echo ">> formatting source"
	@find . -type f -iname '*.go' -not -path './vendor/*' -not -iname '*pb.go' | xargs -L 1 $(GO) fmt

imports:
	@echo ">> fixing source imports"
	@find . -type f -iname '*.go' -not -path './vendor/*' -not -iname '*pb.go' | xargs -L 1 goimports -w

lint:
	@echo ">> linting source"
	@find . -type f -iname '*.go' -not -path './vendor/*' -not -iname '*pb.go' | xargs -L 1 golint

$(DEP_BIN):
	@echo ">> installing package manager(dep)"
	@$(GO) get -u github.com/golang/dep/cmd/dep
