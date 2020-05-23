FROM golang:alpine as builder

RUN adduser -D blog
WORKDIR /
COPY . .

# Compile it
RUN apk add git && \
  go get -u github.com/gorilla/mux && \
  go get github.com/gomarkdown/markdown && \
  go get -u github.com/go-bindata/go-bindata/... && \
  go-bindata -pkg main -o assets.go assets/ && \
  CGO_ENABLED=0 GOOS=linux go build -a \
  -installsuffix cgo \
  -ldflags '-extldflags "-static"' \
  -o blog .

# Create docker
FROM scratch
COPY --from=builder /blog /blog
COPY --from=0 /etc/passwd /etc/passwd
USER blog
WORKDIR /
ENTRYPOINT ["./blog"]
