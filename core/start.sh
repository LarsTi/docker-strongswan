#!/bin/bash

if [[ -z "${PUBLIC_IP}" ]]; then
        echo "Please set env variable \"PUBLIC_IP\" to an IP!"
        exit 1
else
        pubip="${PUBLIC_IP}"
fi

if [[ -z "${MASQUERADE_SUB}" ]]; then
	echo "No Masquerading of Home Network requested"
	echo "This can be a problem on the routing! Take CARE!"
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
if [[ -z "$ip1" || -z "$ip2" ]]; then
	echo "One IP was not found, Error in Configuration"
	echo "This is not recoverable!"
	exit 1
fi 
subnet="$(ip route | grep $me | awk '{print $1}')"
echo "Found subnet $subnet for private ip $me"
if [[ -z "$MASQUERADE_SUB" ]]; then
	echo "Not Masquerading"
else
	echo "Masquerading all Traffic not coming from internal network $subnet, so you can ping it!"
	iptables --table nat --append POSTROUTING ! --source $subnet --jump MASQUERADE
fi
echo "Setting up NAT so strongswan can use the public ip without caring about routing"
iptables --table nat --insert PREROUTING --destination $me --jump DNAT --to-destination $pubip

iptables --table nat --insert POSTROUTING --source $pubip --jump SNAT --to-source $me

echo "Finished Natting"

/usr/libexec/ipsec/charon --debug-ike 1 --debug-cfg 1 --debug-mgr 2
