# docker-strongswan
Strongswan in Docker-Container
## Motivation
I need a reliable IPSec implementation, free to use and free to customize --> https://www.strongswan.org
I want to monitor my IPSec with prometheus --> writing your own exporter
I want to build a webinterface for changing IPSec configurations --> golang
And last but not least i want to learn more about go.
Consider this project a hobby.

## public ips
Have you ever tried to use a public ip address inside your docker container without giving it access to the host network stack?

Have you ever tried to use network namespaces in a docker container?

I stumbled about both questions while investigating how to get things done. Network namespaces are great for a separation and to assign multiple routable IPs to one machine, but docker prohibits the use of network namespaces. So you cant just add a pair of veths in your docker containers and move one of those in a different netns, giving it a public ip.
The idea was clear, adding a public ip to the container (exactly the IP the server has, as we dont want to take any foreign ips), make it routable, and let strognswan do the rest.
Easy idea, hard to build. As mentioned before, netns are not available. For any reason (i could not verify what is the reason), you could just run ip addr add to get your public ip inside the container routable, but strongswan will always go back to NAT and have issues, if the strongswan container is not the initiator (which a server in a roadwarrior scenario is not!).
So the next idea was macvlan, but that did not work either, cause now all the traffic (from the other containers) will be captured by your macvlan container. Something you probably dont want.
Working with a dummy interface and SNAT/DNAT sounds also promising, but i did not manage to get the packets from the dummy interface to eth0.
So i asked around and got an idea: build a network with your public ip via docker-compose and put the ipsec container in that network.
Great idea, but i always got an error "address is already in use", when i tried to build a /32 or /31 subnet. This is caused by docker allocating 2 IPs (if possible, you can create a /32 network, but cant assign a container, as far as i know), one for the gateway and one broadcast address.
So i made a /30 subnet of my public ip (e.g. 4.4.4.6) and gave the container the 4.4.4.6, and the gateway 4.4.4.4/30.
If you want to have 4.4.4.4 on you docker, you would need 4.4.4.0/29 to make 4.4.4.4 a usable ip in this subnet.

## Architecture
I want to use docker, this makes it easy, to adapt everything to a new machine. 
Though docker is very powerful tool, i experienced some issues not easy to solve, due to the architecture of docker.
The Project is (at least for now) separated into two containers:

- ipsec
- mgmt

There is a chance this will increase.

### ipsec
This container is responsible for the ipsec tunneling.
It is connected to two networks:

- private 172.16.0.250/16 --> fix ip (internal ip)
- public 4.4.4.4/29       --> fix ip (public ip)

The public network is very special. Once you try to ping from 4.4.4.4 to 1.1.1.1, you will get an error, due to your host having 4.4.4.4 and answering or rejecting those packages, as it did not know, your container send it. The ping from internal works just fine, as it gets masqueraded. but did you try to masquerade 4.4.4.4 to look as it comes from 4.4.4.4? 
this problem took me a while and i then opted for a SNAT/DNAT:

iptables --table nat --append POSTROUTING --source --$pubip --jump SNAT --to-source $me
iptables --table nat --append PREROUTING --destination $me --jump DNAT --to-destination $pubip

Now you can ping from your pubip and will get a response, and now it is also possible for strongswan to build a connection to anyone in the world, with your public ip, though the container itself has no host access.
To allow the start script to work properly, you have to provide the public ip to use as a environment parameter (PUBLIC_IP).
