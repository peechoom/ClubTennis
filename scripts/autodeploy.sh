#!/usr/bin/env bash

cd /home/alec
if [ -f start_ssh_agent.sh ]; then
	. start_ssh_agent.sh
fi

cd /home/alec/ClubTennis

LOCK_FILE="$(pwd)/deyployment.lock"
rm ./deployment.log

TTL=$( { time flock -n $LOCK_FILE ./scripts/rebuild_if_changed.sh >> ./deployment.log 2>&1; } 2>&1 )

echo "$(date --utc +%FT%TZ): Build completed in $TTL" >> ./deployment_result.log
