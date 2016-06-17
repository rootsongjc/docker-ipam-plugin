# Docker custom network plugin

#### Overview

Create a new docker custom network named "mynet" with a bridge br0.
Docker container allocated on different hosts which started with "mynet" can be assgined one IP address by yourself or automated assgined so that the containers can be accessed by each others.

#### Design

Feature

1.IP manager

All IP addressed stored in ectd. The plugin create a unix socker - /var/run/docker/plugins/talkingdata.sock. Docker invoke the socket when users communcate with docker everytime.

2.Create network

Create a network configuration file on /etc/sysconfig/network-scripts/ifcfg-br0, restart network it works.

For example:

```
DEVICE=br0

TYPE=Bridge

BOOTPROTO=static

IPADDR=10.10.0.85

GATEWAY=10.10.0.1

NETMASK=255.255.192.0

ONBOOT=yes

NOZEROCONF=yes

IPV6INIT=no

NM_CONTROLLED=no

DELAY=0
```

#### Usage
Enviromnent 

OS: CentOS7.1

Gateway: 172.16.16.1

Host network IP range: 172.16.20.54/20 - 172.16.20.58/20
Build the project and put it at /usr/bin/oam-docker-ipam  
Etcd cluster: http://172.16.20.53:2379,http://172.16.20.54:2379,http://172.16.20.55:2379

1.  Assign host IP 
      oam-docker-ipam --cluster-store=http://172.16.20.53:2379，http://172.16.20.54:2379，http://172.16.20.55:2379 host-range --ip-start 172.16.20.54/20 --ip-end 172.16.20.58/20 --gateway 172.16.16.1

2.  Assign docker IP
      oam-docker-ipam --cluster-store=http://172.16.20.53:2379，http://172.16.20.54:2379，http://172.16.20.55:2379 ip-range --ip-start 172.16.21.100/20 --ip-end 172.16.21.200/20

3.  Run docker ipam plugin
      nohup oam-docker-ipam --debug=true --cluster-store=http://172.16.20.53:2379,http://172.16.20.54:2379,http://172.16.20.55:2379 server 2>&1 >> /var/log/oam-docker-ipam.log &

4.  Create custom network
      oam-docker-ipam --cluster-store=http://172.16.20.53:2379，http://172.16.20.54:2379，http://172.16.20.55:2379 create-network --ip 172.16.20.54

####Info

- Usage: oam-docker-ipam --help


- cluster-store defualt http://127.0.0.1:2379


- Open debug mode for more information 


- Docker will choose a random IP address from etcd IP range if you run a docker contianer wihtout assign a specific IP.

#### Reference
https://docs.docker.com/engine/extend/plugins
https://github.com/docker/libnetwork/blob/master/docs/ipam.md
https://github.com/docker/go-plugins-helpers

