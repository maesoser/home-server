# Home Server configuration

## Introduction

I this repo I am going to try to summarize all the configuration that I have on my home server.

## Containers

I've several containers running on my home server. Such container are orchestated using a single docker-compose yaml file. As you can infer from the repo, I use the following foder structure to keep everything organized:

- **config:** In this folder I store the configuration files of each container if needed.
- **data:** In this folder I store the data in use for each container. For instance, the database files on the prometheus container.
- **containers:** If I need to create a container by myself, the Dockerfile and all the needed files are located under this folder.

Here is an schema of the running containers at the time of writting this.

```
              +-------------------+     +-------------------+    +----------+
              |      cotibot      |     |   selenium-hub    |    |  stocks  |
              +-------------------+     +-------------------+    +----------+

              +-------------------+     +-------------------+    +----------+
              |      pihole       |     |  selenium-chrome  |    |   trex   |
              +-------------------+     +-------------------+    +----------+

    9300      +-------------------+
+-------------+  virgin_exporter  |
|             +-------------------+
|
|   9333      +-------------------+
+-------------+  pihole_exporter  |
|             +-------------------+
|
|   9400      +-----------------+
+-------------+  ping_exporter  |
|             +-----------------+
|
|   9134      +-------------------+
+-------------+  docker_exporter  |
|             +-------------------+
|
|   9797      +-------------+    8080
+-------------+  promrelay  +--------------+
|             +-------------+              |
|                                          |
|                                          |
|     +--------------+     +----------+    |
|     |              |9090 |          |3000|
+-----+  prometheus  +-----+ Grafana  +----+
      |              |     |          |    |
      +--------------+     +----------+    |
                                           |
                           +----------+    |
                           |          |    |
                           |   Hugo   +----+
                           |          |    |  +---------+   +---------------+
                           +----------+    |  |         |80 |               |
                                           +--+  nginx  +---+  cloudflared  +->
          +----------+     +----------+    |  |         |   |               |
          |          |3306 |          |80  |  +---------+   +---------------+
          | MariaDB  +-----+  Lychee  +----+
          |          |     |          |    |
          +----------+     +----------+    |
                                           |
                           +----------+    |
                           |          |8080|
                           |  Pihole  +----+
                           |          |    |
                           +----------+    |
                                           |
                          +-----------+    |
                          |           |8080|
                          | httpecho  +----+
                          |           |8443
                          +-----------+

```

### 3th party containers

#### Grafana

Grafana is a famous monitoring solution designed to generate graphics from different data sources like prometheus, elasticsearch, graphite, loki, etc.

#### Prometheus

Prometheus is a well-known TSDB. This service receives data from different local or remote (by using promrelay) exporters.

#### Lychee

[Lychee](https://lychee.electerious.com/) is a self-hosted photo management service, quite easy to configure and maintain.

#### Nginx

I could use something better integrated with Docker like [traefik](https://docs.traefik.io/) but I find nginx easier to configure and maintain and it is very efficient in terms of resources consumption.

Nginx acts as an Edge Router on my configuration, sending the requests to their corresponding container by making use of the [HTTP Host header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Host).

### Self-made containers

#### [ping_exporter](https://github.com/maesoser/home-server/tree/master/containers/ping_exporter)

Ping exporter is an exporter for prometheus that perform TCP or ICMP pings.

#### [docker_exporter](https://github.com/maesoser/home-server/tree/master/containers/docker_exporter)

Exposes docker container metrics like CPU, memory and network usage, among others.

#### [pihole_exporter](https://github.com/maesoser/home-server/tree/master/containers/pihole_exporter)

Exposes metrics about [pihole](https://pi-hole.net/).

#### virgin_exporter

It exposed metrics extracted from Virgin Media Hub router.

#### tplink_exporter

It exposed metrics extracted from tplink router.

#### weather_exporter

Weather exporter is an UDP receiver for a custom weather station based on the ESP8266.

#### [promrelay](https://github.com/maesoser/home-server/tree/master/containers/promrelay)

Promrelay is a relay for prometheus exporters that are on another machine and whose metrics needs to be pushed to a prometheus instance instead of being pushed for such instance.

### Self-made Dockerfiles

#### Cloudflared

[Cloudflared](https://github.com/cloudflare/cloudflared) is a binary created by cloudflare that allow me to expose HTTP services by making use of a technology called [Argo Tunnel](https://www.cloudflare.com/es-es/products/argo-tunnel/). Cloudflare does not offer an official docker image so I created mine.

#### Hugo

#### Jupyter



## Useful functions

```bash
function docker-sh() {
  if [ -z "$1" ]; then
    echo -e "Usage:\n\tdocker-bash <container_name>"
  else
    echo "Logging into container $1"
    docker exec -it $1 /bin/bash
    if [ $? -ne 0 ]; then
      echo "Falling back to /bin/sh"
      docker exec -it $1 /bin/sh
      echo "Exitcode $?"
    fi
  fi
}
```

```bash
function docker-clean() {
  echo -e "Deleting stopped containers"
  docker container prune -f
  echo -e "Deleting unused images"
  docker image prune -af
  # docker rmi $(docker images --filter "dangling=true" -q --no-trunc)
}
```

```bash
function docker-alpine(){
  docker run -it \
    --mount type=bind,source=/media/storage/alpine-data,target=/home \
    --rm alpine /bin/ash
}
```

```bash
function prometheus-reload(){
  curl -X POST http://localhost:9090/-/reload
  echo ""
  echo "TARGETS"
  curl http://localhost:9090/api/v1/targets
  echo ""
  echo "CONFIG Loaded:"
  curl http://localhost:9090/api/v1/status/config
  echo ""
}
```

## Useful aliases

```
export PATH="/home/sierra/.local/bin/:$PATH"
export COMPOSECONF="/media/storage/docker-compose.yml"

alias sys-upgrade='sudo apt update && sudo apt upgrade -y && sudo apt-get autoremove -y'
alias logins='last -a -i -s -7days'
alias onhemera='wakeonlan 00:71:c2:09:ce:37'
alias shutdown="/usr/local/bin/confirm /sbin/shutdown"
alias reboot="/usr/local/bin/confirm /sbin/reboot"
alias freecache='sudo sh -c "echo 1 > /proc/sys/vm/drop_caches"'
alias rotatelogs='sudo logrotate -vf /etc/logrotate.conf'
alias compose-edit="vi $COMPOSECONF"
alias docker-compose="docker-compose -f $COMPOSECONF"
alias sizes="du -h --max-depth=1"
alias dc="docker-compose -f $COMPOSECONF"
alias d="docker"
alias debian='docker run -it --name debian --rm debian /bin/bash'
alias fwlist='sudo iptables -nL --line-numbers'
```

## crontab


```bash
# clean-up and backup
15 0 * * 1 docker image prune -af
30 0 * * 1 nocache rsync -avC --delete --exclude '.local' --exclude '.go' --exclude '.cache' --exclude 'go' /home/user/ /media/backup/home/ > /tmp/home-backup.log 2>&1
45 0 * * 1 nocache rsync -avC --delete --exclude 'docker-data' --exclude 'prometheus/data' /media/storage/ /media/backup/storage/ > /tmp/storage-backup.log 2>&1


# EXPORTERS
2 */6 * * * /usr/local/bin/speedtest_exporter
*/3 * * * * /usr/local/bin/temp_exporter

# DNS UPDATE
*/5 * * * * curl -s -u $EMAIL:$PASSWD "https://now-dns.com/update?hostname=$HOST" > /tmp/now-dns-last-update 2>&1

# Docker restart services at Saturdays
15 0 * * 6 /usr/local/bin/docker-compose -f "/media/storage/infra/docker-compose.yml" restart trex crontab > /dev/null 2>&1
```

## Scripts

#### ddns

```

```
