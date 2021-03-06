#!/bin/bash

iface=$(iw dev | grep Interface | awk '{print $2}')

echo $iface > $SNAP_DATA/interface
sync

# if interface is managed by ifupdown, exit
if [ -e /etc/network/interfaces.d/$iface ]; then
    exit 0
fi

if [ -e $SNAP_DATA/cookie ]; then
    exit 0
fi

#create ap interface
iw dev $iface interface add uap0 type __ap
sleep 2

#Initial wifi interface configuration
ifconfig uap0 up 10.0.60.1 netmask 255.255.255.0
sleep 2

#start dnsmasq
exec $SNAP/usr/sbin/dnsmasq -k -C $SNAP/dnsmasq.conf -l $SNAP_DATA/dnsmasq.leases -x $SNAP_DATA/dnsmasq.pid
