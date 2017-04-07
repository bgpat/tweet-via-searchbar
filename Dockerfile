FROM alpine:3.5

MAINTAINER bgpat

RUN apk --no-cache --update add ca-certificates

COPY bin/tweet-via-searchbar-docker /tweet-via-searchbar
COPY templates /templates

ENTRYPOINT ["/tweet-via-searchbar"]
