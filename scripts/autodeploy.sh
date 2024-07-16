#!/usr/bin/env bash

cd /home/alec
if [ -f start_ssh_agent.sh ]; then
	. start_ssh_agent.sh
fi

cd /home/alec/ClubTennis


# just realized this is a typo, but more than just this script uses it so oh well
LOCK_FILE="$(pwd)/deyployment.lock"

TTL=$( { time flock -n $LOCK_FILE ./scripts/rebuild_if_changed.sh >> ./deployment.log 2>&1; } 2>&1 )
EXIT_STATUS=$?

if [ $EXIT_STATUS -eq 0 ]; then
	echo "$(date --utc +%FT%TZ): Build completed in $TTL" >> ./deployment.log
fi
