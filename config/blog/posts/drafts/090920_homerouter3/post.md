---
title: Adding more capabilities to my router with virtualized firewall
draft: true
tags: networks, projects
---

On my previous post I explained hot to install, configure and use docker to add additional services to my home router. On this post I will explain the tools I developed in order to get a tighter integration between the different elements that make up my router and enhance monitoring.

I found out one really annoying problem. If, for some reason, my router lost power, the firewall VNF that I deployed on it gets corrupted. Although the solution is no so difficult (download the image, uncompress it, boot it up and configure it), it is a little bit annoyin and induces a downtime that could be greatly reduced.

For that reason I started to write a tiny daemon that not only will manage the lifecycle of the firewall and provision it. It will also expose a metrics API and a configuration API and it will include a minimal dashboard to perform basic actions and take a look at the performance of the router.

## routerd

```
GET status/firewall
GET status/containers
GET status/host

GET config/vm
GET config/host
```

### Firewall management

```
{
  status: "BOOTSTRAPED",
  healthy: true
  cpu_pcnt: 34,
  mem_pcnt: 57,
  disk_pcnt: 43,
  disk_read: 1,
  disk_write: 414
  wan_tx: 0,
  wan_rx: 15,
  lan_tx: 124,
  lan_rx: 345,
  image_url: "https://repo.souvlaki.cf/images/opnsense.tar.gz"  
  lan_subnet: "192.168.0.17"
}
```

The different status the VM could have are:
  - UNCONFIGURED: There is no image_url so the vm is undeployed
  - DOWNLOADING: The image_url has been added or modified and the router is downloading the image from internet
  - UNDEPLOYED: The image has been downloaded, it is being copied and it is currently not configured
  - DEPLOYING: The image is being started, configured and then restarted with the right interfaces
  - RUNNING: The image is configured and running
  - STOPPED: The image is configured and it is stopped
  - ERROR: There is an error


### Monitorization


### Webpage



## metricd

Gets statistics about the router

## dockermgrd

Starts/Stops/config containers

## apigwd

communicates the different services with the command an control servers

