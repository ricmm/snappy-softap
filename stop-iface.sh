#!/bin/bash

# remove uap0 interface
iw dev uap0 del

# flush forwarding rules out
iptables --table nat --delete POSTROUTING --out-interface eth0 -j MASQUERADE
iptables --delete FORWARD --in-interface uap0 -j ACCEPT

# disable ipv4 forward
sysctl -w net.ipv4.ip_forward=0
