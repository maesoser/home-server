# PromRelay

PromRelay forwards the metrics gathered from a prometheus exporter to an HTTP endpoint

```
                                       internet
+----------+     +-----------+                          +-----------+       +------------+
|          |     |           |            +  +          |           |       |            |
| exporter +<----+ promrelay +-------->  +  +  +------->+ promrelay +<------+ prometheus |
|          |     |           |  TLS     +  +            |           |       |            |
+----------+     +-----------+                          +-----------+       +------------+

```

Client mode:

```
promrelay 
  --client \
  --compress \ 
  --interval "300s" \ 
  --exporter "http://127.0.0.1:9090/metrics" \
  --exporter "http://127.0.0.1:9300/metrics" \
  --target "metrics.domain.com" \
  --hmac "hmacsecrettext"
```

Server Mode:

```
promrelay
  --server \
  --port 9090 \
  --expose 9191 \
  --noreuse \
  --hmac "hmacsecrettext"
```

To compile it for armv5:

```
env GOOS=linux GOARCH=arm GOARM=5 go build
```
