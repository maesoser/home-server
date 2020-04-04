FROM golang:alpine as builder

WORKDIR $GOPATH/src/github.com/maesoser/glan
COPY . .

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o glan cmd/main.go

FROM alpine

RUN apk add --no-cache nmap nmap-scripts git && \
    git clone https://github.com/vulnersCom/nmap-vulners /usr/share/nmap/scripts/vulners && \
    nmap --script-updatedb

COPY --from=builder /go/src/github.com/maesoser/glan/glan /app/
ENTRYPOINT ["/app/glan", "-targets", "192.168.0.1"]
