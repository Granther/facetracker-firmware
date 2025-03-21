#!/bin/bash

# Update
echo "Updating system"
#apt update && apt upgrade -y;

# General packages
apt install nginx dnsmasq hostapd git neovim wget iptables make -y;

rm /etc/hostapd/hostapd.conf

echo "Pasting in hostapd config"
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

echo "Installing wlan drivers"
git clone https://github.com/RaspAP/raspap-tools
chmod +x raspap-tools/install_wlan_drivers.sh
./raspap-tools/install_wlan_drivers.sh

rm /etc/dnsmasq.conf
rm /etc/nginx/nginx.conf

echo "Pasting in dnsmasq config"
tee -a /etc/dnsmasq.conf << END
interface=uap0
dhcp-range=192.168.4.2,192.168.4.20,255.255.255.0,24h
address=/#/192.168.4.1
END

echo "Stopping all systemd services"
systemctl stop dnsmasq
systemctl stop nignx
systemctl unmask hostapd
systemctl stop hostapd

echo "Enabling IP forwarding"
echo "net.ipv4.ip_forward=1" > /etc/sysctl.conf
sysctl -p

echo "Putting systemd files in correct place"
git clone https://github.com/Granther/facetracker-firmware
cp facetracker-firmware/confs/nginx.conf /etc/nginx/
cp facetracker-firmware/systemd/captive-portal.service /etc/systemd/system/
cp facetracker-firmware/systemd/facetracker.service /etc/systemd/system/
cp facetracker-firmware/systemd/publish-cam.service /etc/systemd/system

systemctl daemon-reload

echo "Installing ARM go"
wget https://go.dev/dl/go1.24.1.linux-arm64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.1.linux-arm64.tar.gz

echo "Adding golang to PATH"
echo "export PATH='$PATH:/usr/local/go/bin'" >> /root/.bashrc
echo "export PATH='$PATH:/usr/local/go/bin'" >> /home/grant/.bashrc

echo "Building captive portal bin"
make -C facetracker-firmware/captive-portal build

echo "Building main facetrack bin"
make -C facetracker-firmware/glorp-track build
