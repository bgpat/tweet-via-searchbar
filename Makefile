NAME           := tweet-via-searchbar
AUTHOR         := bgpat

LDFLAGS        := -ldflags="-s -w"

BIN_FILE  := bin/$(NAME)
SRC_FILES := $(shell find . -name '*.go' -not -path './vendor/*')
VENDOR_FILES := $(shell find . -type 'f' -path './vendor/*')

DOCKER_REGISTRY ?= docker.io
DOCKER_IMAGE_NAME := $(AUTHOR)/$(NAME)
DOCKER_IMAGE_TAG ?= latest
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.DEFAULT_GOAL  := $(BIN_FILE)

$(BIN_FILE): deps $(SRC_FILES)
	go build $(LDFLAGS) -o $@

$(BIN_FILE)-docker: deps $(SRC_FILES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $@

.PHONY: deps
deps: glide glide.yaml $(VENDOR_FILES)
	glide install

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

.PHONY: docker-build
docker-build: $(BIN_FILE)-docker
	docker build -t $(DOCKER_IMAGE) .

.PHONY: ci-build
ci-build: docker-build
	docker push $(DOCKER_IMAGE)

.PHONY: docker-update-latest
docker-update-latest:
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):latest

.PHONY: docker-compose
docker-compose: $(BIN_FILE)-docker
	docker-compose up -d --build

.PHONY: test
test:
	go test -v `glide novendor`

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*
