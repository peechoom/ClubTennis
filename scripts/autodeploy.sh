#!/usr/bin/env bash


LOCK_FILE="$(pwd)/deyployment.lock"
rm ./deployment.log
flock -n $LOCK_FILE ./scripts/rebuild_if_changed.sh >> ./deployment.log 2>&1
