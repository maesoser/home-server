FROM golang:alpine as builder

# Download and install dependencies
RUN apk update \
        && apk upgrade \
        && apk add --no-cache ca-certificates \
        && apk add --no-cache git \
        && update-ca-certificates 2>/dev/null || true

RUN go get gopkg.in/tucnak/telebot.v2 \
	&& go get github.com/dghubble/go-twitter/twitter \
	&& go get github.com/dghubble/oauth1

WORKDIR /app/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o trex .
ENTRYPOINT ["/app/trex"]
