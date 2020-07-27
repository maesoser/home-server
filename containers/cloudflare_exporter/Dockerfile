FROM golang:alpine as builder

LABEL version=1.0

WORKDIR /go/src/cloudflare_exporter
COPY *.go /go/src/cloudflare_exporter/
COPY *.toml /go/src/cloudflare_exporter/

RUN apk update && \
    apk add --no-cache git && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates 2>/dev/null || true && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

ENTRYPOINT ["/go/src/cloudflare_exporter/cloudflare_exporter"]
