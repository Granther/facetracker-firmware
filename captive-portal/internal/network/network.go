package network

// Discover local networks and make permanent connections to them across boot

import (
//	"log"
	"fmt"

//	"github.com/schollz/wifiscan"

	"captive-portal/pkg/utils"
)

func DiscoverNets() {
//        wifis, err := wifiscan.Scan()
//        if err != nil {
//                log.Fatal(err)
//        }
//        for _, w := range wifis {
//               fmt.Println(w.SSID, w.RSSI)
//        }

	// Get all SSID & their corresponding frequencies
	return
}

func ConnectNetwork(ssid, pass string) error {
	return utils.ExecCommand(fmt.Sprintf("raspi-config nonint do_wifi_ssid_passphrase %s %s", ssid, pass))
	// Connect
	// If connection is successful, set state to setup
}
