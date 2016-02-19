#!/bin/bash

# tools are in sbin
export PATH="$PATH:$SNAP_APP_PATH/sbin"

# horrible sleep to make sure everything is up
sleep 10

exec "$SNAP_APP_PATH/bin/softap" $*
