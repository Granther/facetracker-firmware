# Create virtual interface uap0 on boot
iw dev wlan0 interface add uap0 type __ap;
ip addr add 192.168.4.1/24 dev uap0
ip link set dev uap0 up;

systemctl start dnsmasq;
systemctl start hostapd;
systemctl start hostapd;
