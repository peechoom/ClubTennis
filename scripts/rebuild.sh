#!/usr/bin/env bash

git pull

BUILD_VERSION=$(git rev-parse HEAD)

echo "$(date --utc +%FT%TZ): Updating to version $BUILD_VERSION"

docker compose rm -f
docker compose build

OLD_CONTAINER=$(docker ps -aqf "name=server")
echo "$(date --utc +%FT%TZ): Scaling new server up..."
BUILD_VERSION=$BUILD_VERSION docker compose up -d --no-deps --scale server=2 --no-recreate server

sleep 10

echo "$(date --utc +%FT%TZ): Scaling old server down..."
docker container rm -f $OLD_CONTAINER
docker compose up -d --no-deps --scale server=1 --no-recreate server

echo "$(date --utc +%FT%TZ): Reloading Caddy..."
CADDY_CONTAINER=$(docker ps -aqf "name=caddy")
docker exec $CADDY_CONTAINER caddy reload -c /etc/caddy/Caddyfile