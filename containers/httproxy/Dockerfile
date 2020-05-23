FROM golang:alpine as builder

# Copy the code from the host
WORKDIR /
COPY . .

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server .

# Create docker
FROM scratch
COPY --from=builder /server /app/server
RUN adduser -D httproxy
USER httproxy
ENTRYPOINT ["/app/server"]
