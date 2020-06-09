FROM golang:alpine as builder

# Copy the code from the host
WORKDIR $GOPATH/src/github.com/maesoser/promrelay/
COPY . .

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -o promrelay -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Create docker
FROM scratch
COPY --from=builder /go/src/github.com/maesoser/promrelay/promrelay /app/
ENTRYPOINT ["/app/promrelay", "--server", "--noreuse", "--verbose"]
