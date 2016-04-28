#!/bin/bash

if [ -e $SNAP_DATA/cookie ]; then
  exit 0
fi

# wait for uap0 to be setup from dnsmasq
grep uap0 /proc/net/dev &> /dev/null
while [ $? != 0 ]
do
  sleep 5
  grep uap0 /proc/net/dev &> /dev/null
done

#Enable NAT
iptables --flush
iptables --table nat --flush
iptables --delete-chain
iptables --table nat --delete-chain
iptables --table nat --append POSTROUTING --out-interface eth0 -j MASQUERADE
iptables --append FORWARD --in-interface uap0 -j ACCEPT

sysctl -w net.ipv4.ip_forward=1

#copy conf file in place
if [ ! -e $SNAP_DATA/hostapd.conf ] 
then
  cp $SNAP/hostapd.conf $SNAP_DATA/hostapd.conf
fi

#start hostapd
exec $SNAP/usr/sbin/hostapd $SNAP_DATA/hostapd.conf
