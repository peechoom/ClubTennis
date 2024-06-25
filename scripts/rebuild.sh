#!/usr/bin/env bash

git pull

BUILD_VERSION=$(git rev-parse HEAD)

echo "$(date --utc +%FT%TZ): Updating to version $BUILD_VERSION"

docker compose rm -f
docker compose build

OLD_CONTAINER=$(docker ps -aqf "name=server")
echo "$(date --utc +%FT%TZ): Scaling new server up..."
BUILD_VERSION=$BUILD_VERSION docker compose up -d --no-deps --scale server=2 --no-recreate --wait server

NEW_CONTAINER=$(docker ps -aqf "name=server" | grep -v "$OLD_CONTAINER")
SERVER_HEALTHY=$(docker container inspect --format '{{ .State.Health.Status }}' $NEW_CONTAINER)

if [[ "$SERVER_HEALTHY" == "healthy" ]]; then
    echo "$(date --utc +%FT%TZ): Scaling old server down..."
    docker container rm -f $OLD_CONTAINER
else 
    echo "$(date --utc +%FT%TZ): New server integrity could not be verified."
    echo "$(date --utc +%FT%TZ): Reverting changes..."
    docker container rm -f $NEW_CONTAINER
fi

docker compose up -d --no-deps --scale server=1 --no-recreate server

echo "$(date --utc +%FT%TZ): Reloading Caddy..."
CADDY_CONTAINER=$(docker ps -aqf "name=caddy")
docker exec $CADDY_CONTAINER caddy reload -c /etc/caddy/Caddyfile
yes | docker image prune