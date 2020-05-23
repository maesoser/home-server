# Adding SDN capabilities to my home router-firewall


## fwmgrd

Is an utility that manages the firewall vm:

 - Get its status
 - Stop/start the Firewall

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

## metricd

Gets statistics about the router

## dockermgrd

Starts/Stops/config containers

## apigwd

communicates the different services with the command an control servers

