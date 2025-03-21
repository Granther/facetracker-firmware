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

func ConnectNetwork(ssid, pass string, hidden bool) error {
	plain_flag := 1 // Quotes are not present
	hidden_flag := 0 // Not hidden
	if hidden {
		hidden_flag = 1 // Set to hidden
	}

	return utils.ExecCommand(fmt.Sprintf("raspi-config nonint do_wifi_ssid_passphrase %s %s %d %d", ssid, pass, hidden_flag, plain_flag))
	// Connect
	// If connection is successful, set state to setup
}
