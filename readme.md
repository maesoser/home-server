# Home Server configuration

## Introduction

## 3th party containers

## Common containers

## Useful functions

```bash
function docker-bash() {
  if [ -z "$1" ]; then
    echo -e "Usage:\n\tdocker-bash <container_name>"
  else
    echo "Logging into container $1"
    docker exec -it $1 /bin/bash
  fi
}
```

```bash
function docker-clean() {
  docker container prune
  docker image prune -a
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
```

## crontab


```
# clean-up and backup
15 0 * * 1 docker image prune -af
30 0 * * 1 nocache rsync -avC --delete --exclude '.local' --exclude '.go' --exclude '.cache' --exclude 'go' /home/user/ /media/backup/home/ > /tmp/home-backup.log 2>&1
45 0 * * 1 nocache rsync -avC --delete --exclude 'docker-data' --exclude 'prometheus/data' /media/storage/ /media/backup/storage/ > /tmp/storage-backup.log 2>&1


# EXPORTERS
2 */6 * * * /usr/local/bin/speedtest_exporter
*/3 * * * * /usr/local/bin/temp_exporter

# DNS UPDATE
*/5 * * * * curl -s -u $EMAIL:$PASSWD "https://now-dns.com/update?hostname=hecuba.vpndns.net" > /tmp/now-dns-last-update 2>&1

# Docker restart services at Saturdays
15 0 * * 6 /usr/local/bin/docker-compose -f "/media/storage/infra/docker-compose.yml" restart trex crontab > /dev/null 2>&1
```

## Scripts

```

```