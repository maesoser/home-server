FROM golang:alpine as builder

# Copy the code from the host
COPY . /go/src/github.com/maesoser/ping_exporter/
WORKDIR /go/src/github.com/maesoser/ping_exporter


# Compile it
RUN apk add git && go get golang.org/x/net/icmp && \
  go get golang.org/x/net/ipv4 && \
  go get golang.org/x/net/ipv6 && \
  go get github.com/tevino/tcp-shaker
RUN CGO_ENABLED=0 GOOS=linux go build -a \
  -installsuffix cgo \
  -ldflags '-extldflags "-static"' \
  -o /app/ping_exporter cmd/ping_exporter/main.go

# Create docker
FROM scratch
COPY --from=builder /app/ping_exporter /app/
ENTRYPOINT ["/app/ping_exporter", "-port", "9400" ]
