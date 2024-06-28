#!/bin/sh

cd /home/alec/ClubTennis
docker compose up --build > ./server.log 2>&1 &
echo $! | sudo tee /run/docker-webserver.pid > /dev/null
