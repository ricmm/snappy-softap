#!/bin/bash

# tools are in sbin
export PATH="$PATH:$SNAP/sbin"

# horrible sleep to make sure everything is up
sleep 10

exec "$SNAP/bin/softap" $*
