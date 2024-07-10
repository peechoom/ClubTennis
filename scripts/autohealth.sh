#!/usr/bin/env bash

cd /home/alec/ClubTennis

LOCK_FILE="$(pwd)/deyployment.lock"

flock --timeout 3 $LOCK_FILE ./scripts/server_healthcheck.sh >> ./healthcheck.log 2>&1;
