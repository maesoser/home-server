FROM golang:alpine as builder

# Download and install dependencies
RUN apk update && apk add --no-cache git gcc libc-dev && go get github.com/mattn/go-sqlite3

WORKDIR /tmp/stocks/
COPY stocks.go .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o stocks .

# Create docker
FROM scratch
#RUN apk update && apk add ca-certificates
COPY --from=builder /tmp/stocks /stocks/
ENTRYPOINT ["/stocks/stocks", "/stocks/config.json"]
