FROM golang:alpine as builder

# Download and install dependencies
RUN apk update && apk add --no-cache git && go get github.com/oschwald/geoip2-golang && go get github.com/gorilla/mux

# Copy the code from the host
COPY . .

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /maxmind .

# Create docker
FROM scratch
COPY --from=builder /maxmind /maxmind
ENTRYPOINT ["/maxmind"]
