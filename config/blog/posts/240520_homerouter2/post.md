---
title: Installing docker on my x86 router/firewall
draft: false
tags: networks, projects
---


In my [previous post](https://blog.souvlaki.cf/posts/250420_homerouter/) I explained the steps I followed in order to add a virtualized firewall on top of a x86 based router and filter the traffic from the host itself through that firewall.  I had to do that because I wanted to install additional services on my router besides the firewall.

One of the easiest ways to deploy services on a linux machine is by simply spinning up docker containers from [dockerhub](https://hub.docker.com/), so my next obvious step is to install docker and try running some containers.

### Installing docker

Installing docker on the host is easy, we just need to follow the [official guide](https://docs.docker.com/engine/install/debian/) to do it. Besides that I we will also install two other tools that I use a lot together with docker: [ctop](https://github.com/bcicen/ctop) and [docker-compose](https://docs.docker.com/compose/). ctop is a htop like interface that displays container metrics and let you manage them easily from the terminal. docker-compose is a tool that let you define multiple containers on a single yaml file, removing the annoyance of dealing with large and complicated docker commands.

```bash
$ sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
$ curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
$ sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"
$ sudo apt-get update
$ sudo apt-get install docker-ce docker-ce-cli containerd.io
$ sudo curl -L "https://github.com/docker/compose/releases/download/1.25.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
$ sudo wget https://github.com/bcicen/ctop/releases/download/v0.7.3/ctop-0.7.3-linux-amd64 -O /usr/local/bin/ctop
$ sudo chmod +x /usr/local/bin/ctop
```

It is also useful to follow the [extra steps](https://docs.docker.com/engine/install/linux-postinstall/) needed to allow a non privileged user to interact with the docker daemon.

### Debugging Internet access problem

After all these changed I rebooted my router and tragically I came up with one pretty big problem. The PC connected to the router has internet access without any hassle but I wasn't able to get Internet access from the router. I was aware that [docker modifies the iptables config](https://docs.docker.com/network/iptables/) and I suspected that somehow, this was preventing my router to access to internet.

Let's start the debugging session by making a series of pings to see which destinations are unreacheable:

```bash
# Ping to the enp1s0 management interface
$ ping -c 1 10.10.11.1 | grep trans
1 packets transmitted, 1 received, 0% packet loss, time 0ms
# Ping to the veth1 IP address inside LAN bridge
$ ping -c 1 10.10.10.2 | grep trans
1 packets transmitted, 1 received, 0% packet loss, time 0ms
# Ping to vnet0 IP address inside LAN bridge
$ ping -c 1 10.10.10.1 | grep trans
1 packets transmitted, 0 received, 100% packet loss, time 0ms
```

It seems that ICMP traffic is able to arrive at the LAN bridge but it is not being forwarded to the gateway which is `10.10.10.1` .

```bash
$ ip r s
default via 10.10.10.1 dev veth0 
10.10.10.0/24 dev veth0 proto kernel scope link src 10.10.10.2 
10.10.11.0/24 dev enp1s0 proto kernel scope link src 10.10.11.1 
169.254.0.0/16 dev enp1s0 scope link metric 1000 
172.17.0.0/16 dev docker0 proto kernel scope link src 172.17.0.1 linkdown
```

Let's analyze iptables rules and their hits to see if some rule is responsible of this behavior. It is convinient to rememeber that iptables has several tables. Each table is able to perform different actions:

- **NAT:** Performs several address network translation techniques.
- **Raw:** By default, iptables is stateful. This table is used mainly for configuring exemptions from connection tracking. This feature would be useful if our router supported a lot of traffic and we wanted to reduce CPU load on it.
- **Mangle:** This table is used to modify some values inside the IP header of the packets that go through it. It is also used to add marks on the packets at kernel level for further processing in other tables or by other networking tools like [xfrm](http://man7.org/linux/man-pages/man8/ip-xfrm.8.html) (or other utilities included on the iproute2 package).
- **Filter:** This table implements decisions about dropping a packet or letting it go through.

```bash
+------------+                                    +-------------+
|            |                                    |             |
| PREROUTING |                                    | POSTROUTING |
|            |                                    |             |
+------------+                                    +-------------+
| RAW        |                                    | RAW         |
| MANGLE     |                                    | MANGLE      |
| NAT        |                                    | NAT         |
+------------+            +------------+          +-------------+
      |                   |            |             ^     ^
      |                   | FORWARDING |             |     |
  <Routing?>---no-------->+            +-------------+     |
      |                   +------------+                   |
     yes                  | FILTER     |                   |
      |                   | MANGLE     |                   |
      v                   +------------+           +------------+
+------------+                                     |            |
|            |                                     |   OUTPUT   |
|   INPUT    +------------------------------------>|            |
|            |                                     +------------+
+------------+                                     | FILTER     |
| FILTER     |            +------------+           | NAT        |
| MANGLE     +----------->| LOCALHOST  +---------->| MANGLE     |
+------------+            +------------+           | RAW        |
                                                   +------------+
```

There is neither raw nor mangle rules implemented on our router. The configured filtering and NAT rules are the following:

```bash
$ sudo iptables -t nat -L -n -v
Chain PREROUTING (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination         
   14  3876 DOCKER     all  --  *      *       0.0.0.0/0            0.0.0.0/0            ADDRTYPE match dst-type LOCAL        

Chain POSTROUTING (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination         
    0     0 MASQUERADE  all  --  *      !docker0  172.17.0.0/16        0.0.0.0/0           
20983 1497K MASQUERADE  all  --  *      veth0   0.0.0.0/0            0.0.0.0/0           

Chain OUTPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination         
    7   588 DOCKER     all  --  *      *       0.0.0.0/0           !127.0.0.0/8          ADDRTYPE match dst-type LOCAL

Chain DOCKER (2 references)
 pkts bytes target     prot opt in     out     source               destination         
    0     0 RETURN     all  --  docker0 *       0.0.0.0/0            0.0.0.0/0           

$ sudo iptables -t filter -L -n -v
       
Chain FORWARD (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination         
83339   46M DOCKER-USER  all  --  *      *       0.0.0.0/0            0.0.0.0/0           
83339   46M DOCKER-ISOLATION-STAGE-1  all  --  *      *       0.0.0.0/0            0.0.0.0/0           
    0     0 ACCEPT     all  --  *      docker0  0.0.0.0/0            0.0.0.0/0            ctstate RELATED,ESTABLISHED
    0     0 DOCKER     all  --  *      docker0  0.0.0.0/0            0.0.0.0/0           

Chain DOCKER-ISOLATION-STAGE-1 (1 references)
 pkts bytes target     prot opt in     out     source               destination         
    0     0 DOCKER-ISOLATION-STAGE-2  all  --  docker0 !docker0  0.0.0.0/0            0.0.0.0/0           
83339   46M RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0           

Chain DOCKER-ISOLATION-STAGE-2 (1 references)
 pkts bytes target     prot opt in     out     source               destination         
    0     0 DROP       all  --  *      docker0  0.0.0.0/0            0.0.0.0/0           
    0     0 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0           

Chain DOCKER-USER (1 references)
 pkts bytes target     prot opt in     out     source               destination         
83339   46M RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0           
```

There is one [significant difference](https://www.digitalocean.com/community/tutorials/a-deep-dive-into-iptables-and-netfilter-architecture#chain-traversal-order) between doing an ICMP ping from a PC connected to the router (which works) and making the same ping from the router itself. In the first case the packet starts its journey into iptables by hitting prerouting chain. Meanwhile on the second case the packet skips prerouting and forwarding and goes into iptables through output chain.

If we add a log rule into postrouting chain we can see that the ping launched from the router is hitting it.

```bash
sudo iptables -t nat -I POSTROUTING -m limit --limit 2/min -j LOG --log-level 4 --log-prefix 'POSTROUTING '
sudo grep ICMP /var/log/syslog
[48793.474646] POSTROUTING IN= OUT=veth0 SRC=10.10.10.2 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=48081 DF PROTO=ICMP TYPE=8 CODE=0 ID=11775 SEQ=5 
```

If we launch a ping from the PC connected to the router we see the forwarding chain is being hit:

```bash
sudo iptables -t filter -I FORWARD -m limit --limit 2/min -j LOG --log-level 4 --log-prefix 'FORWARD '
sudo grep "ICMP" /var/log/syslog | grep "FORW"
[52757.370311] FORWARD IN=enp1s0 OUT=veth0 MAC=0c:e8:5c:68:32:76:00:e0:4c:a7:ce:3c:08:00 SRC=10.10.11.25 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=63 ID=1682 DF PROTO=ICMP TYPE=8 CODE=0 ID=2328 SEQ=16
```

In the case of the ping sent from the router, it shouldn't hit the NAT rule. Let's try removing it and adding a narrower version that doesn't include the router address:

```bash
$ sudo iptables -L -t nat --line-numbers

Chain POSTROUTING (policy ACCEPT)
num  target     prot opt source               destination         
1    MASQUERADE  all  --  172.17.0.0/16        anywhere            
2    MASQUERADE  all  --  anywhere             anywhere             

$ sudo iptables -t nat -D POSTROUTING 2
$ sudo iptables -t nat -A POSTROUTING ! -s 10.10.10.2 -o veth0 -j MASQUERADE
```

It worked! After modifying the line on `/etc/network/interfaces` our router will have access to internet without any problem.

## Testing docker with a couple of containers

Now let's try adding some containers. The first one will be [Heimdall](https://heimdall.site/), a dashboard that will contain the different access for the services we will add to our router. We will also add [transmission](https://transmissionbt.com/) p2p client and [piHole](https://pi-hole.net/) DNS ad-blocker. Our docker-compose yaml file should look like that:

```yaml
version: "2.4"

services:

  heimdall:
    image: linuxserver/heimdall
    container_name: heimdall
    mem_limit: 128m
    cpu_count: 1
    environment:
      - PUID: '1000'
      - PGID: '1000'
      - TZ: 'Europe/London'
    volumes:
      - ./containers/heimdall:/config
    ports:
      - 80:80
      - 443:443
    restart: unless-stopped

  transmission:
    image: linuxserver/transmission
    container_name: transmission
    mem_limit: 128m
    cpu_count: 1
    environment:
      - PUID: '1000'
      - PGID: '1000'
      - TZ: 'Europe/London'
      - USER: '[redacted]'
      - PASS: '[redacted]'
    volumes:
      - ./containers/transmission:/config
      - ./downloads:/downloads:z
    ports:
      - 9091:9091
      - 51413:51413
      - 51413:51413/udp
    restart: unless-stopped

  pihole:
    container_name: pihole
    image: pihole/pihole:latest
    mem_limit: 128m
    cpu_count: 1
    restart: unless-stopped
    ports:
      - "53:53/tcp"
      - "53:53/udp"
      - "8080:80/tcp"
      - "8443:443/tcp"
    environment:
      TZ: 'Europe/Madrid'
      WEBPASSWORD: '[redacted]'
      VIRTUAL_HOST: '[redacted]'
    volumes:
       - './containers/pihole/etc/:/etc/pihole/'
       - './containers/pihole/dnsmasq/:/etc/dnsmasq.d/'
    cap_add:
       - NET_ADMIN
```

Once we created the needed folders we can bring up the containers by writing:

```bash
docker-compose up -d
```

Docker will download and run the containers. We can check if they're correctly running by using ctop:

![ctop_img](ctop.png)

After configuring the shortcuts on Heimdall we will have a nice dashboard to access to piHole, transmission or our firewall dashboard.

![heimdall dashboard](heimdall.png)

### What next?

Well, there are some things that we could enhance on our router. Maybe a service that will expose CPU/RAM/Temp/Disk/Ifaces metrics as well as a service to redeploy the virtual machine image if it gets corrupted will be nice. Stay tunned!

## References

[A Deep Dive into Iptables and Netfilter Architecture](https://www.digitalocean.com/community/tutorials/a-deep-dive-into-iptables-and-netfilter-architecture)

[IP Masquerading using Iptables](http://billauer.co.il/ipmasq-html.html)

[IP Masquerade and Network Address Translation](https://www.oreilly.com/openbook/linag2/book/ch11.html)

[How to list and delete iptables rules](https://www.digitalocean.com/community/tutorials/how-to-list-and-delete-iptables-firewall-rules)
