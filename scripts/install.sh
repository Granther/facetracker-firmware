#!/bin/bash

# Update
apt update && apt upgrade -y;

# General packages
apt install nginx dnsmasq hostapd git neovim wget -y;

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
rm /etc/nginx/nginx.conf

tee -a /etc/dnsmasq.conf << END
interface=uap0
dhcp-range=192.168.4.2,192.168.4.20,255.255.255.0,24h
address=/#/192.168.4.1
END

systemctl stop dnsmasq
systemctl stop nignx
systemctl unmask hostapd
systemctl stop hostapd

echo "net.ipv4.ip_forward=1" > /etc/sysctl.conf
sysctl -p

#sudo iptables -t nat -A POSTROUTING -o wlan0 -j MASQUERADE
#sudo iptables -A FORWARD -i wlan0 -o uap0 -m state --state RELATED,ESTABLISHED -j ACCEPT
#sudo iptables -A FORWARD -i uap0 -o wlan0 -j ACCEPT

# Create captive portal chain
sudo iptables -t nat -N captiveportal

# Redirect all HTTP traffic to captive portal
sudo iptables -t nat -A captiveportal -p tcp --dport 80 -j DNAT --to-destination 192.168.4.1:80

# Redirect all HTTPS traffic to captive portal
sudo iptables -t nat -A captiveportal -p tcp --dport 443 -j DNAT --to-destination 192.168.4.1:443

# Apply captiveportal chain to all traffic from uap0 except traffic already going to portal
sudo iptables -t nat -A PREROUTING -i uap0 ! -d 192.168.4.1 -j captiveportal

# Force all DNS traffic to your DNS server
sudo iptables -t nat -A PREROUTING -i uap0 -p udp --dport 53 -j DNAT --to 192.168.4.1:53
sudo iptables -t nat -A PREROUTING -i uap0 -p tcp --dport 53 -j DNAT --to 192.168.4.1:53

#iptables -t nat -A PREROUTING -p tcp --dport 80 -j DNAT --to-destination 192.168.4.1

git clone https://github.com/Granther/facetracker-firmware
cp facetracker-firmware/confs/nginx.conf /etc/nginx/
cp facetracker-firmware/systemd/captive-portal.service /etc/systemd/system/
cp facetracker-firmware/systemd/facetracker.service /etc/systemd/system/

systemctl daemon-reload

wget https://go.dev/dl/go1.24.1.linux-arm64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.1.linux-arm64.tar.gz

echo "export PATH='$PATH:/usr/local/go/bin'" > /root/.bashrc
echo "export PATH='$PATH:/usr/local/go/bin'" > /home/grant/.bashrc

