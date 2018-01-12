TARGET              := prom-nginx-exporter
TARGET_SRCS         := $(shell find . -type f -iname '*.go' -not -path './vendor/*')

GO                  := GO15VENDOREXPERIMENT=1 go
DEP_BIN             :=  $(GOPATH)/bin/dep
DOCKER_BUILD_ARGS   ?=
GIT_SUMMARY 		:= $(shell git describe --tags --always 2>/dev/null)

all: test build

fmt:
	@echo ">> checking code style"
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GO) fmt $$d/*.go || ret=$$? ; \
		done ; exit $$ret

vet: vendor
	@echo ">> vetting code"
	@$(GO) vet $(shell $(GO) list ./... | grep -v /vendor/)

build: $(TARGET)

$(TARGET): $(TARGET_SRCS) vendor
	@echo ">> building binary..."
	@echo ">> GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(TARGET)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(TARGET)

docker: GOOS="linux" GOARCH="amd64"
docker: DOCKER_IMAGE_NAME ?= "monitoring-tools/prom-nginx-exporter:$(GIT_SUMMARY)"
docker: Dockerfile build
	@echo ">> building docker image"
	@docker build -t $(DOCKER_IMAGE_NAME) $(DOCKER_BUILD_ARGS) .

vendor: $(DEP_BIN) Gopkg.lock
	@echo ">> installing golang dependencies into vendor directory..."
	@$(DEP_BIN) ensure

$(DEP_BIN):
	@echo "Installing golang dependency manager..."
	@go get -u github.com/golang/dep/cmd/dep

clean:
	@echo ">> cleaning..."
	@rm -rf $(TARGET)
	@rm -rf vendor
