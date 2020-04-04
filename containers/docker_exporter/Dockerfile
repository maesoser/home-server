FROM golang:alpine as builder

WORKDIR /app/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o docker_exporter .

FROM scratch
COPY  --from=builder /app/docker_exporter /app/docker_exporter
ENTRYPOINT ["/app/docker_exporter"]
