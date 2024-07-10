#!/usr/bin/env bash

# get the host from config file (should be ncsutennis.club)
cd /home/alec/ClubTennis


SERVER_CONTAINER=$(docker ps -aqf "name=server")
MYSQL_CONTAINER=$(docker ps -aqf "name=mysql")
CADDY_CONTAINER=$(docker ps -aqf "name=caddy")

SERVER_HEALTHY=$(docker container inspect --format '{{ .State.Health.Status }}' $SERVER_CONTAINER)
MYSQL_HEALTHY=$(docker container inspect --format '{{ .State.Health.Status }}' $MYSQL_CONTAINER)
CADDY_HEALTHY=$(docker container inspect --format '{{ .State.Health.Status }}' $CADDY_CONTAINER)

if [[ "$SERVER_HEALTHY" == "starting" ||  "$MYSQL_HEALTHY" == "starting" || "$CADDY_HEALTHY" == "starting" ]]; then
    # server is still starting, come back later!
    exit 0
fi

if [[ "$SERVER_HEALTHY" != "healthy" ||  "$MYSQL_HEALTHY" != "healthy" || "$CADDY_HEALTHY" != "healthy" ]]; then
    # try to take the compose down and then up again
    docker compose down
    docker compose up -d --wait --wait-timeout 30
else
    exit 0
fi

if [[ "$SERVER_HEALTHY" != "healthy" ||  "$MYSQL_HEALTHY" != "healthy" || "$CADDY_HEALTHY" != "healthy" ]]; then
    # shits hit the fan at this point. Reboot system
    echo "$(date --utc +%FT%TZ): Rebooted server"
    sudo /sbin/reboot
fi

echo "$(date --utc +%FT%TZ): Recovered from slow/no response"
