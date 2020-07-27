# Cloudflare prometheus exporter

This is a fork of the [Cloudflare Exporter](https://gitlab.com/stephane5/cloudflare-prometheus-exporter) written by [Stephane](https://gitlab.com/stephane5).


## Changes

 - Modified the code structure.
 - Removed [cli](https://github.com/urfave/cli) dependency.
 - Moved the collector to an specific structure.
 - You can choose the listening port using `-prom-port`.
 - Now you can add multiple datasets like waf, http or net.
 - If you have several zones but you only want to extract data from one of them you an use `-zone` to specify the name of the zone.
 - The metrics are refreshed when the get request is received by prometheus exporter and not every fixed time, following the guidelines by [Prometheus](https://prometheus.io/docs/instrumenting/writing_exporters/#deployment)
 - Modified "http" dataset to return more metrics:
 - Added "waf" dataset

### Supported metrics

- HTTP
   - Bytes (zoneName)
   - Request (zoneName)
   - CachedBytes (zoneName)
   - CachedRequest (zoneName)
   - EncryptedBytes (zoneName)
   - EncryptedRequests (zoneName)

   - Bytes (contentType, zoneName)
   - Request (contentType, zoneName)

   - Bytes (country, zoneName)
   - Request (country, zoneName)
   - Threats (country, zoneName)

   - Bytes (cacheStatus, contentType, method, zoneName)

   - Requests (sslVersion, zoneName)
   - Requests (HTTPVersion, zoneName)
   - Requests (responseCode, zoneName)

- WAF
   - Events (action, asName, country, ruleID, zoneName)

- Workers
   - CPUTime
   - Errors
   - Requests
   - SubRequests

- Network
   - Bits (attackID)
   - Packets (attackID)

### Format

Here is a sample of metric you should get once running and fetching from the API

```
cloudflare_total_bytes{zoneName="theburritobot.com"} 1.517951e+06
cloudflare_total_requests{zoneName="theburritobot.com"} 67
cloudflare_cached_bytes{zoneName="theburritobot.com"} 865766
cloudflare_cached_requests{zoneName="theburritobot.com"} 15
cloudflare_encrypted_bytes{zoneName="theburritobot.com"} 1.502112e+06
cloudflare_encrypted_requests{zoneName="theburritobot.com"} 66

cloudflare_requests_per_content_type{contentType="txt",zoneName="theburritobot.com"} 3
cloudflare_bytes_per_content_type{contentType="html",zoneName="theburritobot.com"} 824591

cloudflare_requests_per_country{country="US",zoneName="theburritobot.com"} 33
cloudflare_bytes_per_country{country="DK",zoneName="theburritobot.com"} 10514
cloudflare_threats_per_country{country="US",zoneName="theburritobot.com"} 31

cloudflare_processed_bytes{cacheStatus="revalidated",contentType="png",method="GET",zoneName="theburritobot.com"} 69766

cloudflare_requests_per_http_version{version="TLSv1.3",zoneName="theburritobot.com"} 66

cloudflare_requests_per_response_code{responseCode="200",zoneName="theburritobot.com"} 9

cloudflare_requests_per_ssl_type{type="HTTP/1.1",zoneName="theburritobot.com"} 67

cloudflare_waf_events{action="challenge",as="AS-30083-GO-DADDY-COM-LLC",country="US",ruleID="ip",zoneName="theburritobot.com"} 4
```

### Usage

```
cloudflare-prometheus-exporter -h
Usage of ./cloudflare-prometheus-exporter:
  -account string
    	Account ID to be fetched
  -dataset string
    	The data source you want to export, valid values are: http, network (default "http,waf")
  -email string
    	The email address associated with your Cloudflare API token and account
  -key string
    	Your Cloudflare API token
  -prom-port string
    	Prometheus Addr (default "0.0.0.0:2112")
  -zone string
    	Zone Name to be fetched
```

Once launched with valid credentials, the binary will spin a webserver on http://localhost:2112/metrics exposing the metrics received from Cloudflare's GraphQL endpoint.

### Installation

```
go get -u gitlab.com/maesoser/cloudflare-prometheus-exporter
```

Once installed, call it as you would call any other GO binary 

```
cloudflare-prometheus-exporter <options>
```

## TODO

- Add HealthCheck metrics
- Return old last scrapped metrics if time between scrappings is < 5 min
