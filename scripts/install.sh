#!/bin/bash

apt install lighttpd dnsmasq hostapd git neovim -y;

rm /etc/hostapd/hostapd.conf

tee -a /etc/hostapd/hostapd.conf << END
ctrl_interface=/var/run/hostapd
ctrl_interface_group=0
beacon_int=100
auth_algs=1
wpa_key_mgmt=WPA-PSK
ssid=facetrackvr
channel=1
hw_mode=g
wpa_passphrase=TrackMe123
interface=uap0
wpa=2
wpa_pairwise=CCMP
country_code=US
END

git clone https://github.com/RaspAP/raspap-tools
chmod +x raspap-tools/install_wlan_drivers.sh
./raspap-tools/install_wlan_drivers.sh

rm /etc/dnsmasq.conf

tee -a /etc/dnsmasq.conf << END
interface=uap0
dhcp-range=192.168.4.2,192.168.4.20,255.255.255.0,24h
address=/#/192.168.4.1
END

systemctl stop dnsmasq
systemctl unmask hostapd
systemctl stop hostapd

echo "net.ipv4.ip_forward=1" > /etc/sysctl.conf
sysctl -p

sudo iptables -t nat -A POSTROUTING -o wlan0 -j MASQUERADE
sudo iptables -A FORWARD -i wlan0 -o uap0 -m state --state RELATED,ESTABLISHED -j ACCEPT
sudo iptables -A FORWARD -i uap0 -o wlan0 -j ACCEPT

tee -a /var/www/html/index.html << END
<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="refresh" content="0; url=http://example.com" />
</head>
<body>
  <p>If you are not redirected, <a href="http://example.com">click here</a>.</p>
</body>
</html>
END

systemctl stop lighttpd

iptables -t nat -A PREROUTING -p tcp --dport 80 -j DNAT --to-destination 192.168.4.1

