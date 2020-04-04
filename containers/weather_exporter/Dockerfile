FROM golang:alpine as builder

# Copy the code from the host
WORKDIR $GOPATH/src/github.com/maesoser/weather_exporter
COPY . .

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Create docker
FROM scratch
COPY --from=builder /go/src/github.com/maesoser/weather_exporter/weather_exporter /app/
ENTRYPOINT ["/app/weather_exporter"]
