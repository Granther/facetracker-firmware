package setup

import (
	"fmt"
	"time"

	"glorp-track/internal/errors"
	"glorp-track/internal/utils"
	"glorp-track/internal/network"
)

func Init(ifaceName, hotspotIp string) {
	err := setupNetwork(ifaceName, hotspotIp)
	errors.ProcessError(errors.SET_NET_ERR, "unknown error while setting up network for setup mode", err)

	err = setupHotspot()
	errors.ProcessError(errors.SET_HOT_ERR, "unknown error while setting up hotspot for setup mode", err)
}

func Stop() {
	err := deconNetwork()
	errors.ProcessError(errors.SET_NET_ERR, "unknown error while deconstructing network after setup", err)

	err = deconHotspot()
	errors.ProcessError(errors.SET_HOT_ERR, "unknown error while deconstructing hostspot after setup", err)
}

func setupNetwork(ifaceName, hotspotIp string) error {
	if !network.IfaceExists(ifaceName) { // If uap0 does not exist, make it	
		err := utils.ExecCommand(fmt.Sprintf("iw dev wlan0 interface add %s type __ap", ifaceName))
		if err != nil { return err }
	}
	time.Sleep(1 * time.Second) // Sleep becuase this takes time to go into effect

	err := utils.ExecCommand(fmt.Sprintf("ip addr add %s/24 dev %s", hotspotIp, ifaceName))
	if err != nil { return err }

	// Pretty sure this is always 
	err = utils.ExecCommand(fmt.Sprintf("ip link set dev %s up", ifaceName))
	if err != nil { return err }

	return nil
}

func setupHotspot() error {
	err := utils.ExecCommand("systemctl start dnsmasq")
	if err != nil { return fmt.Errorf("unable to start dnsmasq: %w", err) }

	err = utils.ExecCommand("systemctl start captive-portal")
	if err != nil { return fmt.Errorf("unable to start captive-portal: %w", err) }

	err = utils.ExecCommand("systemctl start nginx")
	if err != nil { return fmt.Errorf("unable to start nginx: %w", err) }

	err = utils.ExecCommand("systemctl start hostapd")
	if err != nil { return fmt.Errorf("unable to start hostapd: %w", err) }

	return nil
}

func deconHotspot() error {
	err := utils.ExecCommand("systemctl stop dnsmasq")
	if err != nil { return fmt.Errorf("unable to stop dnsmasq: %w", err) }

	err = utils.ExecCommand("systemctl stop captive-portal")
	if err != nil { return fmt.Errorf("unable to stop captive-portal: %w", err) }

	err = utils.ExecCommand("systemctl stop nginx")
	if err != nil { return fmt.Errorf("unable to stop nginx: %w", err) }

	err = utils.ExecCommand("systemctl stop hostapd")
	if err != nil { return fmt.Errorf("unable to stop hostapd: %w", err) }
	
	return nil
}

func deconNetwork() error {
	/*
	if network.IfaceExists(ifaceName) {
		err := utils.ExecCommand("")
		if err != nil { return fmt.Errorf("unable to remove iface: %s: %w", ifaceName, err) }
	}
	return nil
	*/
	return nil
}

// Create virtual IP, give it IP and make it UP
// Then run all systemd services nessaray
