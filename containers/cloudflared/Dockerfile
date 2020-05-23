FROM golang:alpine as builder

RUN echo 'http://dl-cdn.alpinelinux.org/alpine/edge/main' > /etc/apk/repositories ; \
    echo 'http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories; \
    apk update; \
    apk add --no-cache ca-certificates git gcc build-base; \
    go get -v github.com/cloudflare/cloudflared/cmd/cloudflared

WORKDIR /go/src/github.com/cloudflare/cloudflared/cmd/cloudflared

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' ./

FROM alpine

COPY --from=builder /go/src/github.com/cloudflare/cloudflared/cmd/cloudflared/cloudflared /usr/local/bin/cloudflared
RUN echo 'http://dl-cdn.alpinelinux.org/alpine/edge/main' > /etc/apk/repositories ; \
    echo 'http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories; \
    apk add --no-cache ca-certificates bind-tools; \  
    rm -rf /var/cache/apk/*;

HEALTHCHECK --interval=60s --timeout=5s --start-period=5s CMD nslookup -po=5053 cloudflare.com 127.0.0.1 || exit 1
RUN adduser -D argo
USER argo
ENTRYPOINT ["/usr/local/bin/cloudflared"]
CMD ["tunnel"]
