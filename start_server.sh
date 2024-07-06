#!/bin/sh

cd /home/alec/ClubTennis
if [ -f ./scripts/build_base_image.sh ]; then
	./scripts/build_base_image.sh
fi

docker compose up --build > ./server.log 2>&1 &
echo $! | sudo tee /run/docker-webserver.pid > /dev/null
