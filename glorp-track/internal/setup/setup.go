package setup

import (
	"time"

	"glorp-track/internal/errors"
	"glorp-track/internal/utils"
)

func Init() {
	setupNetwork()
	errors.ProcessError(errors.SET_NET_ERR, "unknown error while setting up network for setup mode", err)

	setupHotspot()
	errors.ProcessError(errors.SET_HOT_ERR, "unknown error while setting up hotspot for setup mode", err)
}

func setupNetwork() error {
	err := utils.ExecCommand("iw dev wlan0 interface add uap0 type __ap")
	if err != nil { return err }
	time.Sleep(1.5 * time.Second) // Sleep becuase this takes time to go into effect

	err = utils.ExecCommand("ip addr add 192.168.4.1/24 dev uap0")
	if err != nil { return err }

	err = utils.ExecCommand("ip link set dev uap0 up")
	if err != nil { return err }

	return nil
}

func setupHotspot() error {
	err := utils.ExecCommand("systemctl start dnsmasq")
	if err != nil { return err }

	err = utils.ExecCommand("systemctl start lighttpd")
	if err != nil { return err }

	err = utils.ExecCommand("systemctl start nginx")
	if err != nil { return err }

	err = utils.ExecCommand("systemctl start hostapd")
	if err != nil { return err }

	return nil
}

// Create virtual IP, give it IP and make it UP
// Then run all systemd services nessaray
