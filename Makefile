NAME := tweet-via-searchbar
BIN_FILE := bin/$(NAME)
SRC_FILES := $(shell find . -name '*.go' -not -path './vendor/*')

.DEFAULT_GOAL := $(BIN_FILE)

$(BIN_FILE): $(SRC_FILES) vendor
	go build -o $@

vendor:
	dep ensure

.PHONY: test
test: $(SRC_FILES)
	go test -v $<

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor
