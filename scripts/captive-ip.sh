#!/bin/bash

# Create captive portal chain
iptables -t nat -N captiveportal

# Redirect all HTTP traffic to captive portal
iptables -t nat -A captiveportal -p tcp --dport 80 -j DNAT --to-destination 192.168.4.1:80

# Redirect all HTTPS traffic to captive portal
iptables -t nat -A captiveportal -p tcp --dport 443 -j DNAT --to-destination 192.168.4.1:443

# Apply captiveportal chain to all traffic from uap0 except traffic already going to portal
iptables -t nat -A PREROUTING -i uap0 ! -d 192.168.4.1 -j captiveportal

# Force all DNS traffic to your DNS server
iptables -t nat -A PREROUTING -i uap0 -p udp --dport 53 -j DNAT --to 192.168.4.1:53
iptables -t nat -A PREROUTING -i uap0 -p tcp --dport 53 -j DNAT --to 192.168.4.1:53

