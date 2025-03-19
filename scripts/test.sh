tee -a ./test.conf << END
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

