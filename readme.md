# Home Server configuration

## Introduction

I this repo I am going to try to summarize all the configuration that I have on my home server.

## Containers

```
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

#### Prometheus

Prometheus is a well-known TSDB. This service receives data from different local or remote (by using promrelay) exporters.

#### Lychee

[Lychee](https://lychee.electerious.com/) is a self-hosted photo management service, quite easy to configure and maintain.

#### Nginx

I could use something better integrated with Docker like [traefik](https://docs.traefik.io/) but I find nginx easier to configure and maintain and it is very efficient in terms of resources consumption.

Nginx acts as an Edge Router on my configuration, sending the requests to their corresponding container by making use of the [HTTP Host header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Host).

### Self-made containers

#### ping_exporter

Ping exporter is an exporter for prometheus that perform TCP or ICMP pings.

#### docker_exporter

Exposes metrics related to docker containers.

#### pihole_exporter

#### virgin_exporter

#### tplink_exporter

#### weather_exporter

#### promrelay



### Self-made Dockerfiles

#### Cloudflared

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
