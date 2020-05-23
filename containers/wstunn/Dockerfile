FROM golang:alpine as builder

RUN adduser -D wstunn
WORKDIR /
COPY wstunn.go .
RUN apk add git && go get github.com/gorilla/websocket
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' wstunn.go

FROM scratch
COPY --from=builder /wstunn /app/
COPY --from=0 /etc/passwd /etc/passwd
USER wstunn
ENTRYPOINT ["/app/wstunn --server 8080"]
