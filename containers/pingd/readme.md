# pingd

PING as a Service

pingd is a daemon that exposes an API to perform pings and retrieve the results

ping request:

```
POST /ping/
{
  host: "www.google.es",
  ipaddr: "",
  type: "tcp/http/icmp/owamp"
  port: 9090
  maxttl: 128,
  is_traceroute: true
  tos : 0
  npackets: -1,
}
```

ping respnose:

```
{
  id: 431234-21341253146-5643632
}
```

result request:

```
GET /ping/{ID}
{
 {
   host: "www.google.es",
   ipaddr: "",
   id: "12341234-1234124-15"
   is_done: false
   type: "tcp/http/icmp/owamp"
   port: 9090
   maxttl: 128,
   is_traceroute: true
   tos : 0
   npackets: -1,
   hops: [{
     from: 1.1.1.1
     to: 2.2.2.2
     fromost: 24124.com
     tohost: 41234.com
     ttl: 12
     sent_num: 12
     recv_num: 12
     loss_pcnt: 0.0
     last_rtt: 2.0
     avg_rtt: 1.3
     std_rtt: 13
     wrst_rtt: 123
     best_rtt: 1231
   }]
 }
}
```

gert

