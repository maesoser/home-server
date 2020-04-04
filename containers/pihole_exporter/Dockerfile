FROM golang:alpine as builder

# Copy the code from the host
WORKDIR /app/
COPY . .

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /app/pihole_exporter .

# Create docker
FROM scratch
COPY --from=builder /app/pihole_exporter /app/
ENTRYPOINT ["/app/pihole_exporter", "-a", "pihole:80" ]
