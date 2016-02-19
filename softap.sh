#!/bin/bash

# tools are in sbin
export PATH="$PATH:$SNAP_APP_PATH/sbin"

exec "$SNAP_APP_PATH/bin/softap" $*
