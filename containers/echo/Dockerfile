FROM golang:alpine as builder

# Copy the code from the host
RUN adduser -D httpecho
WORKDIR /
COPY . .

# Compile it
RUN apk add git && go get -u github.com/gorilla/mux && CGO_ENABLED=0 GOOS=linux go build -a \
  -installsuffix cgo \
  -ldflags '-extldflags "-static"' \
  -o server .

# Create docker
FROM scratch
COPY --from=builder /server /app/server
COPY main.html /app/main.html
COPY --from=0 /etc/passwd /etc/passwd
USER httpecho
ENTRYPOINT ["/app/server"]
