FROM golang:alpine as builder

RUN adduser -D blog
COPY assets /go/src/github.com/blog/assets
COPY pkg /go/src/github.com/blog/pkg
COPY cmd /go/src/github.com/blog/cmd
WORKDIR /go/src/github.com/blog/cmd/

# Compile it
RUN apk add git && \
  go get -u github.com/gorilla/mux && \
  go get github.com/gomarkdown/markdown && \
  go get -u github.com/go-bindata/go-bindata/... && \
  go get github.com/fsnotify/fsnotify && \
  go-bindata -pkg blog \
    -prefix /go/src/github.com/blog/ \
    -o /go/src/github.com/blog/pkg/blog/assets.go /go/src/github.com/blog/assets/ && \
  CGO_ENABLED=0 GOOS=linux go build -a \
  -installsuffix cgo \
  -ldflags '-extldflags "-static"' \
  -o /blog .

# Create docker
FROM scratch
COPY --from=builder /blog /blog
COPY --from=0 /etc/passwd /etc/passwd
USER blog
WORKDIR /
ENTRYPOINT ["./blog"]
