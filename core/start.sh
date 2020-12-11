#!/bin/bash

if [[ -z "${PUBLIC_IP}" ]]; then
        echo "Please set env variable \"PUBLIC_IP\" to an IP!"
        exit 1
else
        pubip="${PUBLIC_IP}"
fi

echo "Expecting 2 IPs (public ip and internal IP)"
ip1="$(hostname -i | awk '{ print $1}' )"
ip2="$(hostname -i | awk '{ print $2}' )"
if [[ "$ip1" == "$pubip" ]]; then
        echo "Found IP1 $ip1 being public IP"
else
        echo "Found IP1 $ip1 being internal IP, routable"
        me="$ip1"
fi
if [[ "$ip2" == "$pubip" ]]; then
        echo "Found IP2 $ip2 being public IP"
else
        echo "Found IP2 $ip2 being internal IP, routable"
        me="$ip2"
fi
 
echo "Setting up NAT so strongswan can use the public ip without caring about routing"
iptables --table nat --insert PREROUTING --destination $me --jump DNAT --to-destination $pubip

iptables --table nat --insert POSTROUTING --source $pubip --jump SNAT --to-source $me

echo "Finished Natting"

/usr/libexec/ipsec/charon --debug-ike 1 --debug-cfg 1 --debug-mgr 2
