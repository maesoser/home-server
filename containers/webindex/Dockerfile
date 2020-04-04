FROM golang:alpine as builder

# Copy the code from the host
WORKDIR /app/
COPY . .

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /app/webserver .

# Create docker
FROM scratch
COPY --from=builder /app/webserver /app/
COPY --from=builder /app/index.html /app/
ENTRYPOINT ["/app/webserver", "--port", "80", "--file", "/app/index.html"]
