#!/bin/bash

iface=$(iw dev | grep Interface | awk '{print $2}')

echo $iface > $SNAP_APP_DATA_PATH/interface

if [ -e $SNAP_DIR/cookie ]; then
    exit 0
fi

#create ap interface
iw dev $iface interface add uap0 type __ap
sleep 2

#Initial wifi interface configuration
ifconfig uap0 up 10.0.60.1 netmask 255.255.255.0
sleep 2

#start dnsmasq
exec $SNAP_APP_PATH/usr/sbin/dnsmasq -k -C $SNAP_APP_PATH/dnsmasq.conf -l $SNAP_APP_DATA_PATH/dnsmasq.leases -x $SNAP_APP_DATA_PATH/dnsmasq.pid
