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

```
promrelay --client --interval "300s" --endpoint "http://127.0.0.1:9090/metrics" --target "relay.souvlaki.cf/system"
promrelay --server --port 9090 --expose 9191 --noreuse
```


```
env GOOS=linux GOARCH=arm GOARM=5 go build
docker-compose -f /media/storage/docker-compose.yml stop promrelay
docker-compose -f /media/storage/docker-compose.yml up -d --build promrelay
./promrelay --client \
  --oneshot --compress \
  --exporter "http://127.0.0.1:9100/metrics" \
  --exporter "http://127.0.0.1:9100/metrics" \
  --exporter "http://127.0.0.1:9100/metrics" \
  --target "https://prom.souvlaki.cf"
sleep 1
docker logs httpecho 2>&1 | tail -n 40
```
