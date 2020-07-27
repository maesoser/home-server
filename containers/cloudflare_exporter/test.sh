

go build
./cloudflare-prometheus-exporter \
  -key "b1fb3e37d3792503ae5ebfd46eaae7cd9d537" \
  -email "massesos@gmail.com" \
  -zone "theburritobot.com" \
  -dataset "waf,http" \
