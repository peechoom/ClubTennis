#!/usr/bin/env bash

# get the host from config file (should be ncsutennis.club)
cd /home/alec/ClubTennis

if [ -f ./config/.env ]; then
    export $(cat ./config/.env | grep SERVER_HOST | xargs)
else
    echo "SERVER_HOST not defined in config/.env"
    exit 1
fi 

# $SERVER_HOST is ur host
URL="https://$SERVER_HOST/ping"
echo pinging $URL
if curl -f --silent -I --max-time 5 --connect-timeout 60 "URL" > /dev/null; then
    exit 0
fi 

# try to take the compose down and then up again
docker compose down
docker compose up -d --wait --wait-timeout 30

CONTAINER=$(docker ps -aqf "name=server")
SERVER_HEALTHY=$(docker container inspect --format '{{ .State.Health.Status }}' $CONTAINER)

if [[ "$SERVER_HEALTHY" != "healthy" ]]; then
    # shits hit the fan at this point. Reboot system
    echo "$(date --utc +%FT%TZ): Rebooted server"
    sudo /sbin/reboot
fi

echo "$(date --utc +%FT%TZ): Recovered from slow/no response"
