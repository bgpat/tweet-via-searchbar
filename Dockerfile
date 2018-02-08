FROM golang:1.9-alpine3.7

RUN apk add -U ca-certificates curl git gcc musl-dev
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 \
		&& chmod +x /usr/local/bin/dep

RUN mkdir -p $GOPATH/src/github.com/bgpat/tweet-via-searchbar
WORKDIR $GOPATH/src/github.com/bgpat/tweet-via-searchbar

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

ADD . ./
RUN go build --ldflags '-s -w -linkmode external -extldflags -static' -o /main


FROM alpine:3.7
RUN apk add -U ca-certificates
COPY --from=0 /main /main
COPY templates /templates
CMD ["/main"]
