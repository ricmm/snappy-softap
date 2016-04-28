#!/bin/bash

if [ ! -e $SNAP_DATA/cookie ]; then
  exit 0
fi

interface=`cat $SNAP_DATA/interface`
cookie=`cat $SNAP_DATA/cookie`

i=0
while [ $i -lt 180 ]; do
	ip=$(ip addr | awk '/inet/ && /'$interface'/{sub(/\/.*$/,"",$2); print $2}')
	echo $ip
	i=$[$i+1]
	curl -i -X POST -H "Content-Type:application/json" http://158.255.238.95/clients/register -d '{"cookie":"'$cookie'", "ip":"'$ip'"}'
	sleep 1
done

rm $SNAP_DATA/cookie
rm $SNAP_DATA/interface
