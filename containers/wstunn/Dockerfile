FROM golang:alpine as builder

RUN adduser -D wstunn
WORKDIR /
COPY wstunn.go .
RUN apk add git && go get github.com/gorilla/websocket
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o wstunn wstunn.go

USER wstunn
WORKDIR /
ENTRYPOINT ["./wstunn", "-buffer", "65534" , "-verbose", "-server", "8080"]
