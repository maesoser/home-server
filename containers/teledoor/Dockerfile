FROM golang:alpine as builder

RUN apk update && apk add --no-cache git && go get github.com/prometheus/client_golang/prometheus && go get gopkg.in/tucnak/telebot.v2

# Copy the code from the host
WORKDIR $GOPATH/src/github.com/maesoser/teledoor
COPY . .

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o teledoor .

# Create docker
FROM scratch
COPY --from=builder /go/src/github.com/maesoser/tplink_exporter/teledoor /app/
ENTRYPOINT ["/app/teledoor"]
